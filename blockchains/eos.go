package bitping

import (
	"fmt"

	types "github.com/auser/bitping/types"
	"github.com/eoscanada/eos-go"
	p2p "github.com/eoscanada/eos-go/p2p"
)

type EosOptions struct {
	Node           string
	NetworkVersion int64
	P2PAddr        string
}

// TODO: make interface for blockchains
type EosApp struct {
	Client  *eos.API
	P2PApi  *p2p.Client
	Info    *eos.InfoResp
	Options EosOptions
}

func NewEosClient(opts EosOptions) (*EosApp, error) {
	fmt.Printf("%#v\n", opts)
	// api := eos.New(opts.Node)

	// info, err := api.GetInfo()
	// if err != nil {
	// 	log.Fatal("Error getting info: %#v\n", err)
	// }

	// cID := info.ChainID
	// p2pApi := p2p.NewClient(opts.P2PAddr, cID, uint16(opts.NetworkVersion))
	app := &EosApp{
		// 	Client:  api,
		// 	P2PApi:  p2pApi,
		// 	Info:    info,
		Options: opts,
	}

	return app, nil
}

type blockHandler struct {
	// p2p.Handler
}

func (b *blockHandler) Handle(msg interface{}) {
	// route := msg.Route
	// fmt.Printf("msg: %#v\n", route)
}

func (app *EosApp) Run(
	blockChan chan types.Block,
	// transChan chan []types.Transaction,
	errChan chan error,
) {

	// client := app.P2PApi
	// info := app.Info

	// b := blockHandler{}

	// client.RegisterHandlerFunc(b.Handle)
	// fmt.Printf("Registered\n")

	// err := client.ConnectAndSync(
	// 	info.HeadBlockNum,
	// 	info.HeadBlockID,
	// 	info.HeadBlockTime.Time,
	// 	0,
	// 	make([]byte, 32))
	// //err = client.ConnectRecent()
	// if err != nil {
	// 	log.Fatal("Error: %s\n", err.Error())
	// }

	// out, err := app.Client.GetBlockByNum(6342084)
	// if err != nil {
	// 	log.Println("get block err")
	// }
	// signedBlock := out.SignedBlock
	// for txNum := range signedBlock.Transactions {
	// 	tx := signedBlock.Transactions[txNum]
	// 	if tx.Transaction.Packed != nil {
	// 		unpacked, err := tx.Transaction.Packed.Unpack()
	// 		if err != nil {
	// 			log.Println(" unpack tractions error  ", err.Error())
	// 		}

	// 		for idx := range unpacked.Actions {
	// 			action := unpacked.Actions[idx]

	// 			if action.Data != nil {
	// 				err := action.MapToRegisteredAction()
	// 				if err != nil {
	// 					log.Fatalf("processAction:ction.MapToRegisteredAction %v", err.Error(), idx)
	// 					continue
	// 				}

	// 				switch op := action.Data.(type) {
	// 				case *token.Transfer:
	// 					log.Printf("Transfer  From : %s , To : %s ", op.From, op.To) // 191
	// 				}

	// 			} else {
	// 				log.Println("no action.data")
	// 			}
	// 		}
	// 	}
	// }

	// info, err := app.Client.GetInfo()
	// if err != nil {
	// 	log.Fatal("Error getting info: %#v\n", err)
	// }

	// done := make(chan bool)
	// cID := info.ChainID

	// client := p2p.NewClient(*app.Options.Node, cID, uint16(*networkVersion))
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
