package avm

import (
	"github.com/icloudland/starchain/crypto"
	. "github.com/icloudland/starchain/errors"
	"github.com/icloudland/starchain/common/log"
	"errors"
	"github.com/icloudland/starchain/common"
)



type ECDsaCrypto struct  {
}

func (c * ECDsaCrypto) Hash160( message []byte ) []byte {
	temp, _ := common.ToCodeHash(message)
	return temp.ToArray()
}

func (c * ECDsaCrypto) Hash256( message []byte ) []byte {
	return []byte{}
}

func (c * ECDsaCrypto) VerifySignature(message []byte,signature []byte, pubkey []byte) (bool,error) {
	var log = log.NewLog()
	log.Debug("message: %x \n", message)
	log.Debug("signature: %x \n", signature)
	log.Debug("pubkey: %x \n", pubkey)

	pk,err := crypto.DecodePoint(pubkey)
	if err != nil {
		return false,NewDetailErr(errors.New("[ECDsaCrypto], crypto.DecodePoint failed."), ErrNoCode, "")
	}

	err = crypto.Verify(*pk, message,signature)
	if err != nil {
		return false,NewDetailErr(errors.New("[ECDsaCrypto], VerifySignature failed."), ErrNoCode, "")
	}

	return true,nil
}
