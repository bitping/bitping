package blockchains

import (
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/auser/bitping/types"
	"github.com/codegangsta/cli"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
)

// EosOptions store the EosApp options
type EosOptions struct {
	Node           string
	NetworkVersion int64
}

// EosApp holds the EOS Client and configuration of an EOS App
// It allows the user to watch for new blockchain blocks generates a Go
// representation of both the original block as well as a unified block
type EosApp struct {
	Client  *eos.API
	Info    *eos.InfoResp
	Options EosOptions
}

// NewEosClient creates a new EosClient
func NewEosClient(opts EosOptions) (*EosApp, error) {
	log.Printf("EOS Opts %v\n", opts)
	api := eos.New(opts.Node)

	info, err := api.GetInfo()
	if err != nil {
		log.Fatalf("GetInfo Error, %v", err)
	}

	app := &EosApp{
		Client:  api,
		Info:    info,
		Options: opts,
	}

	return app, nil
}

// Name returns the app name
func (app EosApp) Name() string {
	return "EOS Watcher"
}

// AddCLIFlags configures the CLI Settings
func (app EosApp) AddCLIFlags(fs []cli.Flag) []cli.Flag {
	return append(fs,
		cli.StringFlag{
			Name:  "eos",
			Usage: "eos address",
		},
		// cli.StringFlag{
		// 	Name:  "eos-p2p",
		// 	Usage: "eos p2p address",
		// 	Value: "peering.mainnet.eoscanada.com:9876",
		// },
		cli.Int64Flag{
			Name:  "eos-version",
			Usage: "eos network version",
			Value: int64(1206),
		},
	)
}

// CanConfigure determines if enough CLI Flags are set to configure the app
func (app EosApp) CanConfigure(c *cli.Context) bool {
	return c.String("eos") != ""
}

// Configure reads CLI Flag settints and configures app
func (app *EosApp) Configure(c *cli.Context) error {
	nodePath := c.String("eos")
	client := eos.New(nodePath)

	info, err := client.GetInfo()
	if err != nil {
		log.Fatalf("GetInfo Error, %v", err)
		return err
	}

	app.Client = client
	app.Info = info
	app.Options = EosOptions{
		Node:           nodePath,
		NetworkVersion: c.Int64("eos-version"),
	}

	return nil
}

// Watch starts running the block watcher
func (app *EosApp) Watch(
	blockCh chan types.Block,
	// transChan chan []types.Transaction,
	errCh chan error,
) {
	log.Printf("Running EOS\n")

	for {
		info := app.Info
		latestInfo, err := app.Client.GetInfo()
		if err != nil {
			log.Fatalf("GetInfo Error: %v", err)
		}
		app.Info = latestInfo

		for blockNum := info.LastIrreversibleBlockNum; blockNum < latestInfo.LastIrreversibleBlockNum; blockNum++ {
			log.Printf("EOS Getting Block: %v", blockNum)
			block, err := app.Client.GetBlockByNum(blockNum) //11819163
			if err != nil {
				log.Fatalf("GetBlockByNum Error: %v", err)
				errCh <- err
				continue
			}

			log.Printf("block: %v", block)

			transactions := make([]types.Transaction, len(block.Transactions))
			for txNum, txReceipt := range block.Transactions {
				log.Printf("tx receipt: %v", txReceipt)

				packedTx := txReceipt.Transaction.Packed

				if packedTx == nil || packedTx.PackedTransaction == nil {
					log.Printf("PackedTransaction has no Transaction")
					continue
				}

				tx, err := packedTx.Unpack()
				if err != nil {
					log.Fatalf("Transaction.Packed.Unpack Error: %v", err)
					errCh <- err
					continue
				}

				statusCode := "unknown"
				switch txReceipt.Status {
				case 0:
					statusCode = "executed"
				case 1:
					statusCode = "soft_fail"
				case 2:
					statusCode = "hard_fail"
				case 3:
					statusCode = "delayed"
				case 4:
					statusCode = "expired"
				case 255:
					statusCode = "unknown"
				}

				trxSigs := make([]string, len(tx.Signatures))
				for i, sig := range tx.Signatures {
					trxSigs[i] = sig.String()
				}

				trxCmp := ""
				switch packedTx.Compression {
				case 0:
					trxCmp = "none"
				case 1:
					trxCmp = "zlib"
				}

				cfd := make([]string, len(tx.ContextFreeData))
				for i, cf := range tx.ContextFreeData {
					cfd[i] = hex.EncodeToString(cf)
				}

				cfActs := make([]types.EOSAction, len(tx.ContextFreeActions))
				for i, cfAct := range tx.ContextFreeActions {
					dat, err := json.Marshal(cfAct.Data)
					if err != nil {
						log.Fatalf("Could not json.Marshal cfAct.Data: %v", err)
						errCh <- err
						continue
					}

					cfActs[i] = types.EOSAction{
						Account: string(cfAct.Account),
						Name:    string(cfAct.Name),
						HexData: hex.EncodeToString(cfAct.HexData),
						Data:    string(dat),
					}
				}

				exts := make([]types.EOSExtension, len(tx.Extensions))
				for i, ext := range tx.Extensions {
					exts[i] = types.EOSExtension{
						Type: uint64(ext.Type),
						Data: hex.EncodeToString(ext.Data),
					}
				}

				id, err := packedTx.ID()
				if err != nil {
					log.Fatalf("PackedTransaction.ID Decode Error: %v", err)
					errCh <- err
					continue
				}

				transactions[txNum] = types.Transaction{
					BlockHash:   hex.EncodeToString(block.ID),
					BlockNumber: int64(block.BlockNum),
					Hash:        hex.EncodeToString([]byte(id)),
					EOSTransactionReceipt: &types.EOSTransactionReceipt{
						Status:               statusCode,
						CPUUsageMicroSeconds: uint64(txReceipt.CPUUsageMicroSeconds),
						NetUsageWords:        uint64(txReceipt.NetUsageWords),
						TRX: types.EOSTransactionWithID{
							ID:                    hex.EncodeToString([]byte(id)),
							Signatures:            trxSigs,
							Compression:           trxCmp,
							PackedTRX:             hex.EncodeToString(packedTx.PackedTransaction),
							PackedContextFreeData: hex.EncodeToString(packedTx.PackedContextFreeData),
							ContextFreeData:       cfd,
							Transaction: types.EOSUnpackedTransaction{
								Expiration:              tx.Expiration.Unix(), // unpacked.Expiration.Unix(),
								RefBlockNum:             uint64(tx.RefBlockNum),
								MaxNetUsageWords:        uint64(tx.MaxCPUUsageMS),
								MaxCPUUsageMicroSeconds: uint64(tx.MaxNetUsageWords),
								DelaySec:                uint64(tx.DelaySec),
								RefBlockPrefix:          uint64(tx.RefBlockPrefix),
								ContextFreeActions:      cfActs,
								TransactionExtensions:   exts,
							},
						},
					},
				}

				acts := make([]types.Action, len(tx.Actions))
				eosActs := make([]types.EOSAction, len(tx.Actions))
				for i, act := range tx.Actions {
					dat, err := json.Marshal(act.Data)
					if err != nil {
						log.Fatalf("Could not json.Marshal act.Data: %v", err)
						errCh <- err
						continue
					}

					eosActs[i] = types.EOSAction{
						Account: string(act.Account),
						Name:    string(act.Name),
						HexData: hex.EncodeToString(act.HexData),
						Data:    string(dat),
					}

					acts[i] = types.Action{
						BlockHash:       hex.EncodeToString(block.ID),
						BlockNumber:     int64(block.BlockNum),
						TransactionHash: transactions[txNum].Hash,
						Address:         eosActs[i].Account,
						Data:            []byte(string(dat)),

						EOSAction: &eosActs[i],
					}

					if act.Data != nil {
						err := act.MapToRegisteredAction()
						// err := act.MapToRegisteredAction()
						if err != nil {
							log.Fatalf("MapToRegisteredAction %v", err)
							errCh <- err
							continue
						}

						switch op := act.Data.(type) {
						case *token.Transfer:
							log.Printf("Transfer From: %s , To: %s, Quantity: %s", op.From, op.To, op.Quantity) // 191
							acts[i].From = string(op.From)
							acts[i].To = string(op.To)
							acts[i].Value = types.BigIntFromInt(op.Quantity.Amount)
							acts[i].Symbol = string(op.Quantity.Symbol.Symbol)
							acts[i].Precision = uint64(op.Quantity.Precision)
						case *token.Create:
							// Send token to self meaning
							log.Printf("Created By: %s, Quantity: %s", op.Issuer, op.MaximumSupply)
							acts[i].From = string(op.Issuer)
							acts[i].To = string(act.Account)
							acts[i].Value = types.BigIntFromInt(op.MaximumSupply.Amount)
							acts[i].Symbol = string(op.MaximumSupply.Symbol.Symbol)
							acts[i].Precision = uint64(op.MaximumSupply.Precision)
						case *token.Issue:
							log.Printf("Created By: %s, Quantity: %s", op.To, op.Quantity)
							acts[i].From = string(act.Account)
							acts[i].To = string(op.To)
							acts[i].Value = types.BigIntFromInt(op.Quantity.Amount)
							acts[i].Symbol = string(op.Quantity.Symbol.Symbol)
							acts[i].Precision = uint64(op.Quantity.Precision)
						}
					} else {
						log.Println("Custom Data")
					}
				}

				transactions[txNum].TRX.Transaction.Actions = eosActs
				transactions[txNum].Actions = acts
			}

			blockObj := types.Block{
				Hash:       hex.EncodeToString(block.ID),
				HeaderHash: hex.EncodeToString(block.ID),
				Network:    "eos",
				Number:     int64(block.BlockNum),
				ParentHash: hex.EncodeToString(block.Previous),
				Time:       block.Timestamp.Unix(),

				EOSBlock: &types.EOSBlock{
					Producer:              string(block.Producer),
					Confirmed:             uint64(block.Confirmed),
					TransactionMerkleRoot: hex.EncodeToString(block.TransactionMRoot),
					ActionMerkleRoot:      hex.EncodeToString(block.ActionMRoot),
					ProducerSignature:     block.ProducerSignature.String(),
					RefBlockPrefix:        uint64(block.RefBlockPrefix),
					ChainID:               hex.EncodeToString(info.ChainID),
				},

				Transactions: transactions,
			}

			blockCh <- blockObj
		}
	}
}
