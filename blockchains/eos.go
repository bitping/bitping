package blockchains

import (
	"encoding/hex"
	"log"

	"github.com/auser/bitping/types"
	"github.com/codegangsta/cli"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
)

type EosOptions struct {
	Node           string
	NetworkVersion int64
}

// TODO: make interface for blockchains
type EosApp struct {
	Client  *eos.API
	Info    *eos.InfoResp
	Options EosOptions
}

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

func (app EosApp) Name() string {
	return "EOS Watcher"
}

func (app EosApp) IsConfigured(c *cli.Context) bool {
	return c.String("eos") != ""
}

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
			singletonTransactions := make([]types.Transaction, 0)
			for txNum, txReceipt := range block.Transactions {
				// packed := tx.Receipt.PackedTransaction

				// if tx.Transaction == nil {
				// 	continue
				// }

				packedTx := txReceipt.Transaction.Packed

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
				case 255:
					statusCode = "unknown"
				}

				// TOOD: Fix this?
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

				cfActs := make([]types.EOSAction, len(tx.Actions))
				for i, cfAct := range tx.ContextFreeActions {
					cfActs[i] = types.EOSAction{
						Account: string(cfAct.Account),
						Name:    string(cfAct.Name),
						HexData: hex.EncodeToString(cfAct.HexData),
						Data:    cfAct.Data,
					}
				}

				exts := make([]types.EOSExtension, len(tx.Extensions))
				for i, ext := range tx.Extensions {
					exts[i] = types.EOSExtension{
						Type: uint64(ext.Type),
						Data: hex.EncodeToString(ext.Data),
					}
				}

				transactions[txNum] = types.Transaction{
					BlockHash:   hex.EncodeToString(block.ID),
					BlockNumber: int64(block.BlockNum),
					Hash:        hex.EncodeToString([]byte(tx.ID())),
					EOSTransactionReceipt: &types.EOSTransactionReceipt{
						Status:               statusCode,
						CPUUsageMicroSeconds: uint64(txReceipt.CPUUsageMicroSeconds),
						// CPUUsageMicroSeconds: uint64(tx.CPUUsageMicroSeconds),
						NetUsageWords: uint64(txReceipt.NetUsageWords),
						TRX: types.EOSTransactionWithID{
							ID:          hex.EncodeToString([]byte(tx.ID())),
							Signatures:  trxSigs,
							Compression: trxCmp,
							// PackedTRX:             hex.EncodeToString(packed.Transaction),
							// PackedContextFreeData: hex.EncodeToString(packed.PackedContextFreeData),
							ContextFreeData: cfd,
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

				acts := make([]types.EOSAction, len(tx.Actions))
				for i, act := range tx.Actions {
					acts[i] = types.EOSAction{
						Account: string(act.Account),
						Name:    string(act.Name),
						HexData: hex.EncodeToString(act.HexData),
						Data:    act.Data,
					}

					if act.Data != nil {
						err := act.MapToRegisteredAction()
						if err != nil {
							log.Fatalf("MapToRegisteredAction %v", err)
							errCh <- err
							continue
						}

						switch op := act.Data.(type) {
						case *token.Transfer:
							log.Printf("Transfer From: %s , To: %s, Quantity: %s", op.From, op.To, op.Quantity) // 191
							transactions[txNum].From = string(op.From)
							transactions[txNum].To = string(op.To)
							transactions[txNum].Value = types.BigIntFromInt(op.Quantity.Amount)
							transactions[txNum].Symbol = string(op.Quantity.Symbol.Symbol)
							transactions[txNum].Precision = uint64(op.Quantity.Precision)
							transactions[txNum].SingletonIndex = len(singletonTransactions)
							singletonTransactions = append(singletonTransactions, transactions[txNum])
						case *token.Create:
							// Send token to self meaning
							log.Printf("Created By: %s, Quantity: %s", op.Issuer, op.MaximumSupply)
							transactions[txNum].From = string(op.Issuer)
							transactions[txNum].To = string(act.Account)
							transactions[txNum].Value = types.BigIntFromInt(op.MaximumSupply.Amount)
							transactions[txNum].Symbol = string(op.MaximumSupply.Symbol.Symbol)
							transactions[txNum].Precision = uint64(op.MaximumSupply.Precision)
							transactions[txNum].SingletonIndex = len(singletonTransactions)
							singletonTransactions = append(singletonTransactions, transactions[txNum])
						case *token.Issue:
							log.Printf("Created By: %s, Quantity: %s", op.To, op.Quantity)
							transactions[txNum].From = string(act.Account)
							transactions[txNum].To = string(op.To)
							transactions[txNum].Value = types.BigIntFromInt(op.Quantity.Amount)
							transactions[txNum].Symbol = string(op.Quantity.Symbol.Symbol)
							transactions[txNum].Precision = uint64(op.Quantity.Precision)
							transactions[txNum].SingletonIndex = len(singletonTransactions)
							singletonTransactions = append(singletonTransactions, transactions[txNum])
						}
					} else {
						log.Println("No Action.Data")
					}
				}

				transactions[txNum].TRX.Transaction.Actions = acts
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

				Transactions:          transactions,
				SingletonTransactions: singletonTransactions,
			}

			blockCh <- blockObj
		}
	}
}
