package net

import (
	"github.com/icloudland/starchain/common"
	"github.com/icloudland/starchain/core/transaction"
	"github.com/icloudland/starchain/events"
	"github.com/icloudland/starchain/crypto"
	"github.com/icloudland/starchain/core/ledger"
	"github.com/icloudland/starchain/net/protocol"
	"github.com/icloudland/starchain/net/node"
	."github.com/icloudland/starchain/errors"
)

type Neter interface {
	GetTxnPool(byCount bool) map[common.Uint256]*transaction.Transaction
	Xmit(interface{}) error
	GetEvent(eventName string) *events.Event
	GetBookKeepersAddrs() ([]*crypto.PubKey, uint64)
	CleanSubmittedTransactions(block *ledger.Block) error
	GetNeighborNoder() []protocol.Noder
	Tx(buf []byte)
	AppendTxnPool(*transaction.Transaction, bool) ErrCode
}

func StartProtocol(pubKey *crypto.PubKey) protocol.Noder {
	net := node.InitNode(pubKey)
	net.ConnectSeeds()

	return net
}
