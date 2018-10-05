package util

import (
	"strconv"
	"errors"
	. "github.com/icloudland/starchain/common"
	"math/rand"
	"github.com/icloudland/starchain/core/transaction"
	"github.com/icloudland/starchain/core/contract"
	"github.com/icloudland/starchain/account"
	"github.com/icloudland/starchain/core/signature"
)

func MakeTransferTransactionExt(assetID Uint256, input []*transaction.UTXOTxInput,
	batchOut []BatchOut, privs map[string]string, addrs []*transaction.TxOutput) (*transaction.Transaction, error) {

	outputNum := len(batchOut)
	if outputNum == 0 {
		return nil, errors.New("nil outputs")
	}

	var expected Fixed64
	output := []*transaction.TxOutput{}
	// construct transaction outputs
	for _, o := range batchOut {
		outputValue, err := StringToFixed64(o.Value)
		if err != nil {
			return nil, err
		}
		expected += outputValue
		address, err := ToScriptHash(o.Address)
		if err != nil {
			return nil, errors.New("invalid address")
		}
		tmp := &transaction.TxOutput{
			AssetID:     assetID,
			Value:       outputValue,
			ProgramHash: address,
		}
		output = append(output, tmp)
	}

	// construct transaction
	tx, err := transaction.NewTransferAssetTransaction(input, output)
	if err != nil {
		return nil, err
	}
	txAttr := transaction.NewTxAttribute(transaction.Nonce, []byte(strconv.FormatInt(rand.Int63(), 10)))
	tx.Attributes = make([]*transaction.TxAttribute, 0)
	tx.Attributes = append(tx.Attributes, &txAttr)

	tx.OutputsRef = addrs

	SignTx(privs, tx)

	return tx, nil
}

func GetAccountByProgramHash(programHash Uint160) *account.Account {
	return nil
}

func SignTx(privs map[string]string, tx *transaction.Transaction) {

	context := contract.NewContractContextExt(tx)

	for _, hash := range context.ProgramHashes {

		//acc := GetAccountByProgramHash(hash)

		addr, _ := hash.ToAddress()
		privateKeyStr := privs[addr]
		a, _ := HexStringToBytes(privateKeyStr)
		acc, _ := account.NewAccountWithPrivatekey(a)

		signdate, err := signature.SignBySigner(tx, acc)
		if err != nil {

		}
		contract, _ := contract.CreateSignatureContract(acc.PublicKey)
		context.AddContract(contract, acc.PublicKey, signdate)

	}

	tx.SetPrograms(context.GetPrograms())
}
