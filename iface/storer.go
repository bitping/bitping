package iface

// Storage interface is for each storage solution
type Storer interface {
	// Has to be configured
	Configurable

	Name() string

	Push(interface{}) bool
}
