package dbft

import (
	"sync"
	"github.com/icloudland/starchain/crypto"
	."github.com/icloudland/starchain/common"
	tx"github.com/icloudland/starchain/core/transaction"
	"github.com/icloudland/starchain/core/ledger"
	"github.com/icloudland/starchain/common/log"
	msg"github.com/icloudland/starchain/net/message"
	"fmt"
	"github.com/icloudland/starchain/common/serialization"
	cl"github.com/icloudland/starchain/account"

	"github.com/icloudland/starchain/net"
)

const ContextVersion uint32 = 0

type ConsensusContext struct {
	State           ConsensusState
	PrevHash        Uint256
	Height          uint32
	ViewNumber      byte
	BookKeepers     []*crypto.PubKey
	NextBookKeepers []*crypto.PubKey
	Owner           *crypto.PubKey
	BookKeeperIndex int
	PrimaryIndex    uint32
	Timestamp       uint32
	Nonce           uint64
	NextBookKeeper  Uint160
	Transactions    []*tx.Transaction
	Signatures      [][]byte
	ExpectedView    []byte

	header *ledger.Block

	contextMu           sync.Mutex
	isBookKeeperChanged bool
	nmChangedblkHeight  uint32
}


func (cxt *ConsensusContext) M() int {
	return len(cxt.BookKeepers) - (len(cxt.BookKeepers)-1)/3
}

func (cxt *ConsensusContext) ChangeView(viewNum byte) {
	p := (cxt.Height - uint32(viewNum)) % uint32(len(cxt.BookKeepers))
	cxt.State &= SignatureSent
	cxt.ViewNumber = viewNum
	if p >= 0 {
		cxt.PrimaryIndex = uint32(p)
	} else {
		cxt.PrimaryIndex = uint32(p) + uint32(len(cxt.BookKeepers))
	}

	if cxt.State == Initial {
		cxt.Transactions = nil
		cxt.Signatures = make([][]byte, len(cxt.BookKeepers))
		cxt.header = nil
	}
}

func (cxt *ConsensusContext) MakeChangeView() *msg.ConsensusPayload {
	cv := &ChangeView{
		NewViewNumber: cxt.ExpectedView[cxt.BookKeeperIndex],
	}
	cv.msgData.Type = ChangeViewMsg
	return cxt.MakePayload(cv)
}

func (cxt *ConsensusContext) MakeHeader() *ledger.Block {
	if cxt.Transactions == nil {
		return nil
	}
	if cxt.header == nil {
		txHash := []Uint256{}
		for _, t := range cxt.Transactions {
			txHash = append(txHash, t.Hash())
		}
		txRoot, err := crypto.ComputeRoot(txHash)
		if err != nil {
			return nil
		}
		blockData := &ledger.Blockdata{
			Version:          ContextVersion,
			PrevBlockHash:    cxt.PrevHash,
			TransactionsRoot: txRoot,
			Timestamp:        cxt.Timestamp,
			Height:           cxt.Height,
			ConsensusData:    cxt.Nonce,
			NextBookKeeper:   cxt.NextBookKeeper,
		}
		cxt.header = &ledger.Block{
			Blockdata:    blockData,
			Transactions: []*tx.Transaction{},
		}
	}
	return cxt.header
}

func (cxt *ConsensusContext) MakePayload(message ConsensusMessage) *msg.ConsensusPayload {
	message.ConsensusMessageData().ViewNumber = cxt.ViewNumber
	return &msg.ConsensusPayload{
		Version:         ContextVersion,
		PrevHash:        cxt.PrevHash,
		Height:          cxt.Height,
		BookKeeperIndex: uint16(cxt.BookKeeperIndex),
		Timestamp:       cxt.Timestamp,
		Data:            serialization.ToArray(message),
		Owner:           cxt.Owner,
	}
}

func (cxt *ConsensusContext) MakePrepareRequest() *msg.ConsensusPayload {
	preReq := &PrepareRequest{
		Nonce:          cxt.Nonce,
		NextBookKeeper: cxt.NextBookKeeper,
		Transactions:   cxt.Transactions,
		Signature:      cxt.Signatures[cxt.BookKeeperIndex],
	}
	preReq.msgData.Type = PrepareRequestMsg
	return cxt.MakePayload(preReq)
}

func (cxt *ConsensusContext) MakePrepareResponse(signature []byte) *msg.ConsensusPayload {
	preRes := &PrepareResponse{
		Signature: signature,
	}
	preRes.msgData.Type = PrepareResponseMsg
	return cxt.MakePayload(preRes)
}

func (cxt *ConsensusContext) GetSignaturesCount() (count int) {
	count = 0
	for _, sig := range cxt.Signatures {
		if sig != nil {
			count += 1
		}
	}
	return count
}

func (cxt *ConsensusContext) GetStateDetail() string {

	return fmt.Sprintf("Initial: %t, Primary: %t, Backup: %t, RequestSent: %t, RequestReceived: %t, SignatureSent: %t, BlockGenerated: %t, ",
		cxt.State.HasFlag(Initial),
		cxt.State.HasFlag(Primary),
		cxt.State.HasFlag(Backup),
		cxt.State.HasFlag(RequestSent),
		cxt.State.HasFlag(RequestReceived),
		cxt.State.HasFlag(SignatureSent),
		cxt.State.HasFlag(BlockGenerated))

}

func (cxt *ConsensusContext) Reset(client cl.Client, localNode net.Neter) {
	var log = log.NewLog()
	cxt.State = Initial
	cxt.PrevHash = ledger.DefaultLedger.Blockchain.CurrentBlockHash()
	cxt.Height = ledger.DefaultLedger.Blockchain.BlockHeight + 1
	cxt.ViewNumber = 0
	cxt.BookKeeperIndex = -1

	cxt.BookKeepers, cxt.NextBookKeepers, _ = ledger.DefaultLedger.Store.GetBookKeeperList()
	log.Info("curr bookkeeper, len:", len(cxt.BookKeepers))
	log.Info("next bookkeeper, len:", len(cxt.NextBookKeepers))

	var err error
	cxt.NextBookKeeper, err = ledger.GetBookKeeperAddress(cxt.NextBookKeepers)
	if err != nil {
		log.Error("[ConsensusContext] GetBookKeeperAddres failed")
	}

	cxt.Owner = cxt.BookKeepers[0]

	bookKeeperLen := len(cxt.BookKeepers)
	cxt.PrimaryIndex = cxt.Height % uint32(bookKeeperLen)
	cxt.Transactions = nil
	cxt.header = nil
	cxt.Signatures = make([][]byte, bookKeeperLen)
	cxt.ExpectedView = make([]byte, bookKeeperLen)

	for i := 0; i < bookKeeperLen; i++ {
		ac, _ := client.GetDefaultAccount()
		if ac.PublicKey.X.Cmp(cxt.BookKeepers[i].X) == 0 {
			cxt.BookKeeperIndex = i
			cxt.Owner = cxt.BookKeepers[i]
			break
		}
	}

}
