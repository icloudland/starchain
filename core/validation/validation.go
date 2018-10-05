package validation

import (
	sig"github.com/icloudland/starchain/core/signature"
	"errors"
	"github.com/icloudland/starchain/vm/avm/interfaces"
	"github.com/icloudland/starchain/vm/avm"
	."github.com/icloudland/starchain/common"
	."github.com/icloudland/starchain/errors"
	"github.com/icloudland/starchain/crypto"
	"github.com/icloudland/starchain/common/log"
)

func VerifySignableData(signableData sig.SignableData)(bool,error){
	var log = log.NewLog()
	hashes,err := signableData.GetProgramHashes()
	if err != nil {
		log.Tracef("cat't get programhash")
		return false,err
	}
	programs := signableData.GetPrograms()

	if len(hashes) != len(programs){
		return false,errors.New("the data and program have different length")
	}
	for i := 0; i < len(programs); i++ {
		temp, _ := ToCodeHash(programs[i].Code)
		if hashes[i] != temp {
			return false, errors.New("The data hashes is different with corresponding program code.")
		}
		//execute program on VM
		var cryptos interfaces.ICrypto
		cryptos = new(avm.ECDsaCrypto)
		se := avm.NewExecutionEngine(signableData, cryptos, nil, nil, Fixed64(0))
		se.LoadCode(programs[i].Code, false)
		se.LoadCode(programs[i].Parameter, true)
		err := se.Execute()

		if err != nil {
			return false, NewDetailErr(err, ErrNoCode, "")
		}

		if se.GetState() != avm.HALT {
			return false, NewDetailErr(errors.New("[VM] Finish State not equal to HALT."), ErrNoCode, "")
		}

		if se.GetEvaluationStack().Count() != 1 {
			return false, NewDetailErr(errors.New("[VM] Execute Engine Stack Count Error."), ErrNoCode, "")
		}

		flag := se.GetExecuteResult()
		if !flag {
			return false, NewDetailErr(errors.New("[VM] Check Sig FALSE."), ErrNoCode, "")
		}
	}

	return true, nil
}

func VerifySignature(signableData sig.SignableData,pubkey *crypto.PubKey,signature []byte)(bool,error){
	err := crypto.Verify(*pubkey,sig.GetHashData(signableData),signature)
	if err != nil {
		return false,NewDetailErr(err,ErrNoCode,"verfysingature failed")
	}
	return true,nil
}
