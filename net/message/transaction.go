package message

import (
	"github.com/icloudland/starchain/common"
	"github.com/icloudland/starchain/core/transaction"
	"github.com/icloudland/starchain/common/log"
	"github.com/icloudland/starchain/net/protocol"
	"errors"
	."github.com/icloudland/starchain/errors"
	"bytes"
	"encoding/binary"
	"github.com/icloudland/starchain/core/ledger"
	"crypto/sha256"
)

type dataReq struct {
	msgHdr
	dataType common.InventoryType
	hash     common.Uint256
}

type trn struct {
	msgHdr
	// TBD
	//txn []byte
	txn transaction.Transaction
	//hash common.Uint256
}


func (msg trn) Handle(node protocol.Noder) error {
	var log = log.NewLog()
	log.Debug("RX Transaction message")
	tx := &msg.txn
	if !node.LocalNode().ExistedID(tx.Hash()) {
		if errCode := node.LocalNode().AppendTxnPool(&(msg.txn), true); errCode != ErrNoError {
			return errors.New("[message] VerifyTransaction failed when AppendTxnPool.")
		}
		//broadcast tx to neighour node
		node.LocalNode().Relay(node, tx)
		log.Info("Relay transaction")
		node.LocalNode().IncRxTxnCnt()
		log.Debug("RX Transaction message hash", msg.txn.Hash())
		log.Debug("RX Transaction message type", msg.txn.TxType)
	}

	return nil
}

func reqTxnData(node protocol.Noder, hash common.Uint256) error {
	var msg dataReq
	msg.dataType = common.TRANSACTION
	// TODO handle the hash array case
	//msg.hash = hash

	buf, _ := msg.Serialization()
	go node.Tx(buf)
	return nil
}

func (msg dataReq) Serialization() ([]byte, error) {
	hdrBuf, err := msg.msgHdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	err = binary.Write(buf, binary.LittleEndian, msg.dataType)
	if err != nil {
		return nil, err
	}
	msg.hash.Serialize(buf)

	return buf.Bytes(), err
}

func (msg *dataReq) Deserialization(p []byte) error {
	var log = log.NewLog()
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(msg.msgHdr))
	if err != nil {
		log.Warn("Parse datareq message hdr error")
		return errors.New("Parse datareq message hdr error")
	}

	err = binary.Read(buf, binary.LittleEndian, &(msg.dataType))
	if err != nil {
		log.Warn("Parse datareq message dataType error")
		return errors.New("Parse datareq message dataType error")
	}

	err = msg.hash.Deserialize(buf)
	if err != nil {
		log.Warn("Parse datareq message hash error")
		return errors.New("Parse datareq message hash error")
	}
	return nil
}


func NewTxnFromHash(hash common.Uint256) (*transaction.Transaction, error) {
	var log = log.NewLog()
	txn, err := ledger.DefaultLedger.GetTransactionWithHash(hash)
	if err != nil {
		log.Error("Get transaction with hash error: ", err.Error())
		return nil, err
	}

	return txn, nil
}

func NewTxn(txn *transaction.Transaction) ([]byte, error) {
	var log = log.NewLog()
	var msg trn

	msg.msgHdr.Magic = protocol.NETMAGIC
	cmd := "tx"
	copy(msg.msgHdr.CMD[0:len(cmd)], cmd)
	tmpBuffer := bytes.NewBuffer([]byte{})
	txn.Serialize(tmpBuffer)
	msg.txn = *txn
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, tmpBuffer.Bytes())
	if err != nil {
		log.Error("Binary Write failed at new Msg")
		return nil, err
	}
	s := sha256.Sum256(b.Bytes())
	s2 := s[:]
	s = sha256.Sum256(s2)
	buf := bytes.NewBuffer(s[:4])
	binary.Read(buf, binary.LittleEndian, &(msg.msgHdr.Checksum))
	msg.msgHdr.Length = uint32(len(b.Bytes()))
	log.Debug("The message payload length is ", msg.msgHdr.Length)

	m, err := msg.Serialization()
	if err != nil {
		log.Error("Error Convert net message ", err.Error())
		return nil, err
	}

	return m, nil
}

func (msg trn) Serialization() ([]byte, error) {
	hdrBuf, err := msg.msgHdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	msg.txn.Serialize(buf)

	return buf.Bytes(), err
}

func (msg *trn) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(msg.msgHdr))
	err = msg.txn.Deserialize(buf)
	if err != nil {
		return err
	}

	return nil
}

type txnPool struct {
	msgHdr
	//TBD
}

func ReqTxnPool(node protocol.Noder) error {
	msg := AllocMsg("txnpool", 0)
	buf, _ := msg.Serialization()
	go node.Tx(buf)

	return nil
}

