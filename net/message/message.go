package message

import (
	."github.com/icloudland/starchain/net/protocol"
	"bytes"
	"github.com/icloudland/starchain/common/log"
	"errors"
	"encoding/binary"
	"fmt"
	"encoding/hex"
	"crypto/sha256"
)

type Messager interface {
	Verify([]byte) error
	Serialization() ([]byte, error)
	Deserialization([]byte) error
	Handle(Noder) error
}

type msgHdr struct {
	Magic uint32
	CMD 	[MSGCMDLEN]byte
	Length 	uint32
	Checksum	[CHECKSUMLEN]byte
}

type msgCont struct {
	hdr msgHdr
	p interface{}
}

type varStr struct {
	len uint
	buf []byte
}

type filteradd struct {
	msgHdr
}
type filterclear struct {
	msgHdr
}
type filterload struct {
	msgHdr
}

func AllocMsg(t string, length int) Messager {
	var log = log.NewLog()
	//log.Info("get msg:"+t)
	switch t {
	case "msgheader":
		return &msgHdr{}
	case "version":
		var msg version
		// TODO fill the header and type
		copy(msg.Hdr.CMD[0:len(t)], t)
		return &msg
	case "verack":
		var msg verACK
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "getheaders":
		var msg headersReq
		// TODO fill the header and type
		copy(msg.hdr.CMD[0:len(t)], t)
		return &msg
	case "headers":
		var msg blkHeader
		copy(msg.hdr.CMD[0:len(t)], t)
		return &msg
	case "getaddr":
		var msg addrReq
		copy(msg.Hdr.CMD[0:len(t)], t)
		return &msg
	case "addr":
		var msg addr
		copy(msg.hdr.CMD[0:len(t)], t)
		return &msg
	case "inv":
		var msg Inv
		copy(msg.Hdr.CMD[0:len(t)], t)
		// the 1 is the inv type lenght
		msg.P.Blk = make([]byte, length-MSGHDRLEN-1)
		return &msg
	case "getdata":
		var msg dataReq
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "block":
		var msg block
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "tx":
		var msg trn
		copy(msg.msgHdr.CMD[0:len(t)], t)
		//if (message.Payload.Length <= 1024 * 1024)
		//OnInventoryReceived(Transaction.DeserializeFrom(message.Payload));
		return &msg
	case "consensus":
		var msg consensus
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "filteradd":
		var msg filteradd
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "filterclear":
		var msg filterclear
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "filterload":
		var msg filterload
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "getblocks":
		var msg blocksReq
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "txnpool":
		var msg txnPool
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "alert":
		log.Warn("Not supported message type - alert")
		return nil
	case "merkleblock":
		log.Warn("Not supported message type - merkleblock")
		return nil
	case "notfound":
		var msg notFound
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "ping":
		var msg ping
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "pong":
		var msg pong
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "reject":
		log.Warn("Not supported message type - reject")
		return nil
	default:
		log.Warn("Unknown message type")
		return nil
	}
}

func MsgType(buf []byte) (string, error) {
	cmd := buf[CMDOFFSET : CMDOFFSET+MSGCMDLEN]
	n := bytes.IndexByte(cmd, 0)
	if n < 0 || n >= MSGCMDLEN {
		return "", errors.New("Unexpected length of CMD command")
	}
	s := string(cmd[:n])
	return s, nil
}


func NewMsg(t string, n Noder) ([]byte, error) {
	switch t {
	case "version":
		return NewVersion(n)
	case "verack":
		return NewVerack()
	case "getheaders":
		return NewHeadersReq()
	case "getaddr":
		return newGetAddr()

	default:
		return nil, errors.New("Unknown message type")
	}
}

func HandleNodeMsg(node Noder, buf []byte, len int) error {
	var log = log.NewLog()
	if len < MSGHDRLEN {
		log.Warn("Unexpected size of received message")
		return errors.New("Unexpected size of received message")
	}

	log.Debugf("Received data len:  %d\n%x", len, buf[:len])

	s, err := MsgType(buf)
	//log.Info("msg type:",s)
	if err != nil {
		log.Error("Message type parsing error")
		return err
	}

	msg := AllocMsg(s, len)
	if msg == nil {
		log.Error(fmt.Sprintf("Allocation message %s failed", s))
		return errors.New("Allocation message failed")
	}
	// Todo attach a node pointer to each message
	// Todo drop the message when verify/deseria packet error
	msg.Deserialization(buf[:len])
	err = msg.Verify(buf[MSGHDRLEN:len])
	if err != nil {
		log.Error("msg is not invalid")
	}

	return msg.Handle(node)
}


func magicVerify(magic uint32) bool {
	if magic != NETMAGIC {
		return false
	}
	return true
}

func ValidMsgHdr(buf []byte) bool {
	var h msgHdr
	h.Deserialization(buf)
	//TODO: verify hdr checksum
	return magicVerify(h.Magic)
}

func PayloadLen(buf []byte) (int) {
	var h msgHdr
	h.Deserialization(buf)
	return int(h.Length)
}

func LocateMsgHdr(buf []byte) []byte {
	var h msgHdr
	for i := 0; i <= len(buf)-MSGHDRLEN; i++ {
		if magicVerify(binary.LittleEndian.Uint32(buf[i:])) {
			buf = append(buf[:0], buf[i:]...)
			h.Deserialization(buf)
			return buf
		}
	}
	return nil
}

func checkSum(p []byte) []byte {
	t := sha256.Sum256(p)
	s := sha256.Sum256(t[:])

	// Currently we only need the front 4 bytes as checksum
	return s[:CHECKSUMLEN]
}

func reverse(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	return append(reverse(input[1:]), input[0])
}

func (hdr *msgHdr) init(cmd string, checksum []byte, length uint32) {
	hdr.Magic = NETMAGIC
	copy(hdr.CMD[0:uint32(len(cmd))], cmd)
	copy(hdr.Checksum[:], checksum[:CHECKSUMLEN])
	hdr.Length = length
	//hdr.ID = id
}

// Verify the message header information
// @p payload of the message
func (hdr msgHdr) Verify(buf []byte) error {
	var log = log.NewLog()
	if magicVerify(hdr.Magic) == false {
		log.Warn(fmt.Sprintf("Unmatched magic number 0x%0x", hdr.Magic))
		return errors.New("Unmatched magic number")
	}
	checkSum := checkSum(buf)
	if bytes.Equal(hdr.Checksum[:], checkSum[:]) == false {
		str1 := hex.EncodeToString(hdr.Checksum[:])
		str2 := hex.EncodeToString(checkSum[:])
		log.Warn(fmt.Sprintf("Message Checksum error, Received checksum %s Wanted checksum: %s",
			str1, str2))
		return errors.New("Message Checksum error")
	}

	return nil
}



func (msg *msgHdr) Deserialization(p []byte) error {

	buf := bytes.NewBuffer(p[0:MSGHDRLEN])
	err := binary.Read(buf, binary.LittleEndian, msg)
	return err
}

// FIXME how to avoid duplicate serial/deserial function as
// most of them are the same
func (hdr msgHdr) Serialization() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, hdr)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (hdr msgHdr) Handle(n Noder) error {
	// TBD
	return nil
}

