package rpchttp

import (
	. "github.com/icloudland/starchain/account"
	. "github.com/icloudland/starchain/common"
	"github.com/icloudland/starchain/common/log"
	. "github.com/icloudland/starchain/core/asset"
	"github.com/icloudland/starchain/core/contract"
	"github.com/icloudland/starchain/core/signature"
	"github.com/icloudland/starchain/core/transaction"
	"strconv"
)

const (
	ASSETPREFIX = "STC"
)

func NewRegTx(rand string, index int, admin, issuer *Account) *transaction.Transaction {
	name := ASSETPREFIX + "-" + strconv.Itoa(index) + "-" + rand
	description := "description"
	asset := &Asset{name, description, byte(MaxPrecision), AssetType(Share), UTXO}
	amount := Fixed64(1000)
	controller, _ := contract.CreateSignatureContract(admin.PubKey())
	tx, _ := transaction.NewRegisterAssetTransaction(asset, amount, issuer.PubKey(), controller.ProgramHash)
	return tx
}

func SignTx(admin *Account, tx *transaction.Transaction) {
	var log = log.NewLog()
	signdate, err := signature.SignBySigner(tx, admin)
	if err != nil {
		log.Error(err, "signdate SignBySigner failed")
	}
	transactionContract, _ := contract.CreateSignatureContract(admin.PublicKey)
	transactionContractContext := contract.NewContractContext(tx)
	transactionContractContext.AddContract(transactionContract, admin.PublicKey, signdate)
	tx.SetPrograms(transactionContractContext.GetPrograms())
}
