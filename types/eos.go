package types

// https://github.com/EOSIO/eosjs/blob/master/src/schema/chain_types.json
// {
//   "timestamp": "2018-08-23T17:39:23.500",
//   "producer": "acryptolions",
//   "confirmed": 0,
//   "previous": "00aab1dbad726a64d7a034073c63135e04ef8673bd5e9ef09e91d164bf2935ed",
//   "transaction_mroot": "c9e3c0fa1d236d0cd5ec6527b3f54e954cf6e76595640f1bfd7e8705418ea643",
//   "action_mroot": "207a06bd344056e456ab84ca4f271ca14f1617e18dde9060110f4290159d2d9c",
//   "schedule_version": 196,
//   "new_producers": null,
//   "header_extensions": [],
//   "producer_signature": "SIG_K1_KcCBD54ZyX1JiQqpjMDuN99XZi9kR8wPprqUSgHwUwfdHuMrYauNSh47bb9gHmTZeBz5YLc8qYMWHywBt4qNE5sqUNYLHC",
//   "id": "00aab1dce0ee6cf0457088907ac017f32356dc34790b9a03a4f5147dbce3b3b8",
//   "block_num": 11186652,
//   "ref_block_prefix": 2424860741,
//   "transactions": [],
// }

type EOSBlock struct {
	Producer              string `json:"producer"`
	Confirmed             uint64 `json:"confirmed"`
	TransactionMerkleRoot string `json:"transactionMroot"`
	ActionMerkleRoot      string `json:"actionMroot"`
	ScheduleVersion       uint64 `json:"scheduleVersion"`
	// NewProducers       [] `json:"newProducers"`
	// HeaderExtensions   [] `json:"headerExtensions"`
	ProducerSignature string `json:"producerSignature"`
	RefBlockPrefix    uint64 `json:"refBlockPrefix"`
	ChainID           string `json:"chainID"`
}

// {
//       "status": "executed",
//       "cpu_usage_us": 1696,
//       "net_usage_words": 28,
//       "trx": {
//         "id": "4528f0f169031bd02402ef096095b4801cb6766b7462f8cfc24626cb2dbd8e92",
//         "signatures": [
//           "SIG_K1_KerL5zixYTk7ojcosM2iG8v2vpu1bbocQXoipmQNevjvRXGxhT6AcND9Ts7MTd1rFMsmn6W2LT3BqZqUYss853YfuhtJUf"
//         ],
//         "compression": "none",
//         "packed_context_free_data": "",
//         "context_free_data": [],
//         "packed_trx": "c0f17e5b8cb0b810b9890000000001808ec958e5ab983b000080f1a86c52d501808ec958e5ab983b00000000a8ed32328101808ec958e5ab983b366ade6765010000d0b1aa000000000017323031382d30382d32335431373a33393a31372e3530302a090000000000004030306161623164303562636364373364323830323736346232633935366165613736343736353065336230336463363238306134393534646462353161326530010000000000000000",
//         "transaction": {
//           "expiration": "2018-08-23T17:41:20",
//           "ref_block_num": 45196,
//           "ref_block_prefix": 2310607032,
//           "max_net_usage_words": 0,
//           "max_cpu_usage_ms": 0,
//           "delay_sec": 0,
//           "context_free_actions": [],
//           "actions": [
//             {
//               "account": "bigertestabc",
//               "name": "updateblk",
//               "authorization": [
//                 {
//                   "actor": "bigertestabc",
//                   "permission": "active"
//                 }
//               ],
//               "data": {
//                 "account": "bigertestabc",
//                 "dealId": "1535045954102",
//                 "eosRealblock": 11186640,
//                 "blktime": "2018-08-23T17:39:17.500",
//                 "blockSize": 2346,
//                 "blockId": "00aab1d05bccd73d2802764b2c956aea7647650e3b03dc6280a4954ddb51a2e0",
//                 "txCounts": 1
//               },
//               "hex_data": "808ec958e5ab983b366ade6765010000d0b1aa000000000017323031382d30382d32335431373a33393a31372e3530302a0900000000000040303061616231643035626363643733643238303237363462326339353661656137363437363530653362303364633632383061343935346464623531613265300100000000000000"
//             }
//           ],
//           "transaction_extensions": []
//         }
//       }
//     }

type EOSExtension struct {
	Type uint64 `json:"type"`
	Data string `json:"data"`
}

type EOSPermissionLevel struct {
	Actor       string `json:"actor"`
	Permisssion string `json:"permission"`
}

type EOSAction struct {
	Account       string               `json:"account"`
	Name          string               `json:"name"`
	Authorization []EOSPermissionLevel `json:"authorization"`
	HexData       string               `json:"hexData"`
	Data          string               `json:"data"`
}

type EOSUnpackedTransaction struct {
	Expiration              int64          `json:"expiration"`
	RefBlockNum             uint64         `json:"refBlockNum"`
	RefBlockPrefix          uint64         `json:"refBlockPrefix"`
	MaxNetUsageWords        uint64         `json:"maxNetUsageWords"`
	MaxCPUUsageMicroSeconds uint64         `json:"maxCPUUsageMS"`
	DelaySec                uint64         `json:"delaySec"`
	Actions                 []EOSAction    `json:"actions"`
	ContextFreeActions      []EOSAction    `json:"contextFreeActions"`
	TransactionExtensions   []EOSExtension `json:"transactionExtensions"`
}

type EOSTransactionWithID struct {
	ID                    string                 `json:"id"`
	Signatures            []string               `json:"signatures"`
	Compression           string                 `json:"compression"`
	PackedTRX             string                 `json:"packedTRX"`
	PackedContextFreeData string                 `json:"packedContextFreeData"`
	ContextFreeData       []string               `json:"contextFreeData"`
	Transaction           EOSUnpackedTransaction `json:"transaction"`
}

type EOSTransactionReceipt struct {
	Status               string               `json:"status"`
	CPUUsageMicroSeconds uint64               `json:"cpuUsageUS"`
	NetUsageWords        uint64               `json:"netUsageWords"`
	TRX                  EOSTransactionWithID `json:"trx"`
}
