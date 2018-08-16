package work

// TODO: add persistence
// based largely on this, but with a bolt backend for persistence
// https://github.com/gammazero/workerpool/blob/master/workerpool.go

import (
	"sync"
	"time"

	"github.com/gammazero/deque"
)

const (
	// Value is the size of the queue that workers register their availability on
	readyQueueSize = 16
	// worker pool receives no work for period of time, then stop a goroutine
	idleTimeoutSec = 5
)

// New Creates and starts a pool of worker goroutines
func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		taskQueue:    make(chan func(), 1),
		maxWorkers:   maxWorkers,
		readyWorkers: make(chan chan func(), readyQueueSize),
		timeout:      time.Second * idleTimeoutSec,
		stoppedChan:  make(chan struct{}),
	}

	// start the task dispatcher
	go pool.dispatch()

	return pool
}

// WorkerPool is the main pool structure
type WorkerPool struct {
	maxWorkers   int
	timeout      time.Duration
	taskQueue    chan func()
	readyWorkers chan chan func()
	stoppedChan  chan struct{}
	waitingQueue deque.Deque
	stopMutex    sync.Mutex
	stopped      bool
}

// Stop stops the worker pool and waits for only currently running tasks
// to be complete. Pending tasks that are not currently running are left for dead
// and tasks musst not be submitted after the pool has called stop.
func (p *WorkerPool) Stop() {
	p.stop(false)
}

// StopWait stops the workerpool and waits for all queued tasks to complete
func (p *WorkerPool) StopWait() {
	p.stop(true)
}

// Stopped returns true if this worker pool has been stopped
func (p *WorkerPool) Stopped() bool {
	p.stopMutex.Lock()
	defer p.stopMutex.Unlock()
	return p.stopped
}

// Submit enqueues a function for a worker to execute
//
// Any external values needed by the task function must be captured
// in a closure. Any return values should be returned over a channel that is
// captured in the task function closure.
//
// Submit will not block regardless of the number of tasks submitted. Each task
// is immediately given to an available worker or passed to a goroutine to be given
// to the next available worker. If there are no available workers, the dispatcher
// adds a worker, until the maximum number of workers has been reached is running.
//
// After the max number of workers has been reached and no workers are available,
// incoming tasks are put into a queue and will be executed as soon as workers
// become available
//
// When no new tasks are submitted for a time period and a worker is available, the
// worker is shutdown. As long as no new tasks arrive, one available worker is shutdown each time period has passed until there are no more idle workers. Since the time
// to start new goroutines is insignificant, there is no need to retain idle workers.
func (p *WorkerPool) Submit(task func()) {
	if task != nil {
		p.taskQueue <- task
	}
}

// SubmitWait enqueues the given task and waits for its execution
func (p *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}

	doneChan := make(chan struct{})
	p.taskQueue <- func() {
		task()
		close(doneChan)
	}
	<-doneChan
}

// dispatch sends the next queued task to an available worker
func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan)
	timeout := time.NewTimer(p.timeout)
	var (
		workerCount    int
		task           func()
		ok, wait       bool
		workerTaskChan chan func()
	)

	startReady := make(chan chan func())
Loop:
	for {
		// as long as tasks are in the waiting queue, remove and
		// execute these tasks as workers become available and plaace
		// new incoming tasks on the queue. Once the queue is empty, then
		// go back to submitting incoming tasks directly to available
		// workers
		if p.waitingQueue.Len() != 0 {
			select {
			case task, ok = <-p.taskQueue:
				if !ok {
					break Loop
				}
				if task == nil {
					wait = true
					break Loop
				}
				p.waitingQueue.PushBack(task)
			case workerTaskChan = <-p.readyWorkers:
				// a worker is ready
				workerTaskChan <- p.waitingQueue.PopFront().(func())
			}
			continue
		}
		timeout.Reset(p.timeout)
		select {
		case task, ok = <-p.taskQueue:
			if !ok || task == nil {
				break Loop
			}
			// got some work
			select {
			case workerTaskChan = <-p.readyWorkers:
				workerTaskChan <- task
			default:
				// no worker is ready
				// create a new worker if not at max
				if workerCount < p.maxWorkers {
					workerCount++
					go func(t func()) {
						startWorker(startReady, p.readyWorkers)
						taskChan := <-startReady
						taskChan <- t
					}(task)
				} else {
					p.waitingQueue.PushBack(task)
				}
			}
		case <-timeout.C:
			// Timed out waiting for work to arrive
			// so let's kill a waiting worker
			if workerCount > 0 {
				select {
				case workerTaskChan = <-p.readyWorkers:
					close(workerTaskChan)
					workerCount--
				default:
					// no work, but no ready workers either
					// all workers are busy
				}
			}
		}
	}

	// If instructe to wait for all queued tasks, then remove from the queue
	// and give to workers until queue is empty
	if wait {
		for p.waitingQueue.Len() != 0 {
			workerTaskChan = <-p.readyWorkers
			// a worker is ready
			workerTaskChan <- p.waitingQueue.PopFront().(func())
		}
	}

	// stop all remaining workers as they become ready
	for workerCount > 0 {
		workerTaskChan = <-p.readyWorkers
		close(workerTaskChan)
		workerCount--
	}
}

// startWorker starts a goroutine that executes tasks given by dispatcher
//
// when a new worker starts up, it registers its availability on the startReady
// channel. This ensures that the goroutine associated with the starting worker
// gets to use the worker to execute its task. Otherwise, the main dispatcher loop could steal the new worker and not know how to start up another
// worker for the waiting goroutine. The task would then have to wait for another
// existing worker to become available, even though capacity is
// availablle to start additional workers.
//
// A worker registers its available to do work by putting its task channel on the
// readyWorkers channel, and then writes a task to the worker over the
// worker's task channel. To stop a worker, the dispatcher closes a worker's
// task channel, instead of writing a task to it.
func startWorker(startReady, readyWorkers chan chan func()) {
	go func() {
		taskChan := make(chan func())
		var task func()
		var ok bool
		// register availability
		startReady <- taskChan
		for {
			// read task from dispatcher
			task, ok = <-taskChan
			if !ok {
				// dispatcher told worker to stop
				break
			}

			// Execute the task
			task()

			// register availability
			readyWorkers <- taskChan
		}
	}()
}

// stop tells the dispatcher to exit and if it should complete queued tasks
func (p *WorkerPool) stop(wait bool) {
	p.stopMutex.Lock()
	defer p.stopMutex.Unlock()
	if p.stopped {
		return
	}
	p.stopped = true
	if wait {
		p.taskQueue <- nil
	}

	// close task queue and wait for currently running tasks to complete
	close(p.taskQueue)
	<-p.stoppedChan
}
