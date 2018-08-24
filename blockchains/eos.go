package blockchains

import (
	"fmt"
	"log"

	types "github.com/auser/bitping/types"
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
	fmt.Printf("%#v\n", opts)
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

func (app *EosApp) Run(
	blockCh chan types.Block,
	// transChan chan []types.Transaction,
	errCh chan error,
) {
	fmt.Printf("Running EOS\n")

	for {
		info := app.Info
		latestInfo, err := app.Client.GetInfo()
		if err != nil {
			log.Fatalf("GetInfo Error: %v", err)
		}
		app.Info = latestInfo

		for blockNum := info.LastIrreversibleBlockNum; blockNum < latestInfo.LastIrreversibleBlockNum; blockNum++ {
			block, err := app.Client.GetBlockByNum(blockNum)
			if err != nil {
				log.Fatalf("GetBlockByNum Error: %v", err)
				errCh <- err
				continue
			}

			transactions := make([]types.Transaction, len(block.Transactions))
			singletonTransactions := make([]types.Transaction, len(block.Transactions))
			for txNum, tx := range block.Transactions {
				if tx.Transaction.Packed == nil {
					continue
				}

				packed := tx.Transaction.Packed
				unpacked, err := packed.Unpack()
				if err != nil {
					log.Fatalf("Transaction.Packed.Unpack Error: %v", err)
					errCh <- err
					continue
				}

				statusCode := ""
				switch tx.Status {
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

				trxSigs := make([]string, len(unpacked.Signatures))
				for i, sig := range unpacked.Signatures {
					trxSigs[i] = sig.String()
				}

				trxCmp := ""
				switch packed.Compression {
				case 0:
					trxCmp = "none"
				case 1:
					trxCmp = "zlib"
				}

				cfd := make([]string, len(unpacked.ContextFreeData))
				for i, cf := range unpacked.ContextFreeData {
					cfd[i] = string(cf)
				}

				transactions[txNum] = types.Transaction{
					BlockHash:   string(block.ID),
					BlockNumber: int64(block.BlockNum),
					Hash:        string(tx.Transaction.ID),
					EOSTransactionReceipt: &types.EOSTransactionReceipt{
						Status:               statusCode,
						CPUUsageMicroSeconds: uint64(tx.CPUUsageMicroSeconds),
						NetUsageWords:        uint64(tx.NetUsageWords),
						TRX: types.EOSTransactionWithID{
							ID:                    string(tx.Transaction.ID),
							Signatures:            trxSigs,
							Compression:           trxCmp,
							PackedTRX:             string(packed.PackedTransaction),
							PackedContextFreeData: string(packed.PackedContextFreeData),
							ContextFreeData:       cfd,
							Transaction: types.EOSUnpackedTransaction{
								Expiration:              unpacked.Expiration.Unix(),
								RefBlockNum:             uint64(unpacked.RefBlockNum),
								RefBlockPrefix:          uint64(unpacked.RefBlockPrefix),
								MaxNetUsageWords:        uint64(unpacked.MaxNetUsageWords),
								MaxCPUUsageMicroSeconds: uint64(unpacked.MaxCPUUsageMS),
								DelaySec:                uint64(unpacked.DelaySec),
								// ContextFreeActions:      cfas,
							},
						},
					},
				}

				cfActs := make([]types.EOSAction, len(unpacked.Actions))
				for i, cfAct := range unpacked.ContextFreeActions {
					cfActs[i] = types.EOSAction{
						Account: string(cfAct.Account),
						Name:    string(cfAct.Name),
						HexData: string(cfAct.HexData),
						Data:    cfAct.Data,
					}
				}
				transactions[txNum].TRX.Transaction.ContextFreeActions = cfActs

				exts := make([]types.EOSExtension, len(unpacked.Extensions))
				for i, ext := range unpacked.Extensions {
					exts[i] = types.EOSExtension{
						Type: uint64(ext.Type),
						Data: string(ext.Data),
					}
				}
				transactions[txNum].TRX.Transaction.ContextFreeActions = cfActs

				acts := make([]types.EOSAction, len(unpacked.Actions))
				for i, act := range unpacked.Actions {
					acts[i] = types.EOSAction{
						Account: string(act.Account),
						Name:    string(act.Name),
						HexData: string(act.HexData),
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
							singletonTransactions = append(singletonTransactions, transactions[txNum])
						case *token.Create:
							// Send token to self meaning
							log.Printf("Created By: %s, Quantity: %s", op.Issuer, op.MaximumSupply)
							transactions[txNum].From = string(op.Issuer)
							transactions[txNum].To = string(act.Account)
							transactions[txNum].Value = types.BigIntFromInt(op.MaximumSupply.Amount)
							transactions[txNum].Symbol = string(op.MaximumSupply.Symbol.Symbol)
							transactions[txNum].Precision = uint64(op.MaximumSupply.Precision)
							singletonTransactions = append(singletonTransactions, transactions[txNum])
						case *token.Issue:
							log.Printf("Created By: %s, Quantity: %s", op.To, op.Quantity)
							transactions[txNum].From = string(act.Account)
							transactions[txNum].To = string(op.To)
							transactions[txNum].Value = types.BigIntFromInt(op.Quantity.Amount)
							transactions[txNum].Symbol = string(op.Quantity.Symbol.Symbol)
							transactions[txNum].Precision = uint64(op.Quantity.Precision)
							singletonTransactions = append(singletonTransactions, transactions[txNum])
						}

					} else {
						log.Println("no action.data")
					}
				}

				transactions[txNum].TRX.Transaction.Actions = acts
			}

			blockObj := types.Block{
				Hash:       string(block.ID),
				HeaderHash: string(block.ID),
				Network:    "eos",
				Number:     int64(block.BlockNum),
				ParentHash: string(block.Previous),
				Time:       block.Timestamp.Unix(),

				EOSBlock: &types.EOSBlock{
					Producer:              string(block.Producer),
					Confirmed:             uint64(block.Confirmed),
					TransactionMerkleRoot: string(block.TransactionMRoot),
					ActionMerkleRoot:      string(block.ActionMRoot),
					ProducerSignature:     block.ProducerSignature.String(),
					RefBlockPrefix:        uint64(block.RefBlockPrefix),
					ChainID:               string(info.ChainID),
				},

				Transactions:          transactions,
				SingletonTransactions: singletonTransactions,
			}

			blockCh <- blockObj
		}
	}
}
