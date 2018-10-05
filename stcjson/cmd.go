package stcjson

// GetBlockCountCmd defines the getblockcount JSON-RPC command.
type GetBlockCountCmd struct{}

// NewGetBlockCountCmd returns a new instance which can be used to issue a
// getblockcount JSON-RPC command.
func NewGetBlockCountCmd() *GetBlockCountCmd {
	return &GetBlockCountCmd{}
}

// GetBlockCountCmd defines the getblockcount JSON-RPC command.
type GetBlockCmd struct {
	BlockHeight int64
}

// NewGetBlockCmd returns a new instance which can be used to issue a
// getblock JSON-RPC command.
func NewGetBlockCmd(blockHeight int64) *GetBlockCmd {
	return &GetBlockCmd{
		BlockHeight: blockHeight,
	}
}

type GetRawTransactionCmd struct {
	TxHash string
}

func NewGetRawTransactionCmd(txHash string) *GetRawTransactionCmd {
	return &GetRawTransactionCmd{
		TxHash: txHash,
	}
}

type SendRawTransactionCmd struct {
	RawTx string
}

func NewSendRawTransactionCmd(rawTx string) *SendRawTransactionCmd {
	return &SendRawTransactionCmd{
		RawTx: rawTx,
	}
}

func init() {
	// No special flags for commands in this file.
	flags := UsageFlag(0)

	MustRegisterCmd("getblockcount", (*GetBlockCountCmd)(nil), flags)
	MustRegisterCmd("getblock", (*GetBlockCmd)(nil), flags)
	MustRegisterCmd("getrawtransaction", (*GetRawTransactionCmd)(nil), flags)
	MustRegisterCmd("sendrawtransaction", (*SendRawTransactionCmd)(nil), flags)
}
