package main

import (
	stclog"github.com/icloudland/starchain/common/log"
	"github.com/icloudland/starchain/account"
	"github.com/icloudland/starchain/crypto"
	"os"
	"github.com/icloudland/starchain/core/transaction"
	"github.com/icloudland/starchain/common/config"
	"github.com/icloudland/starchain/core/store/ChainStore"
	"github.com/icloudland/starchain/net"
	"github.com/icloudland/starchain/net/rpchttp"
	"time"
	"github.com/icloudland/starchain/consensus/dbft"
	"github.com/icloudland/starchain/net/httprestful"
	"github.com/icloudland/starchain/core/ledger"
	"github.com/icloudland/starchain/net/protocol"
)

var log = stclog.NewLog()

func main(){
	var err error
	ledger.DefaultLedger = new(ledger.Ledger)
	ledger.DefaultLedger.Store,err = ChainStore.NewLedgerStore()
	defer ledger.DefaultLedger.Store.Close()
	if err != nil {
		log.Fatal("open LedgerStore err:", err)
		os.Exit(1)
	}
	//init store
	//ledger.DefaultLedger.Store.InitLedgerStore(ledger.DefaultLedger)
	transaction.TxStore = ledger.DefaultLedger.Store
	crypto.SetAlg(config.Parameters.EncryptAlg)
	ledger.StandbyBookKeepers = account.GetBookKeepers()
	//create gesesis block if the first time start program
	chain, err := ledger.GenesisBlock(ledger.StandbyBookKeepers)
	checkErr(err,"generate blockchain failed")
	ledger.DefaultLedger.Blockchain = chain
	log.Info("open the wallet")
	cli := account.GetClient()
	if cli == nil {
		log.Fatal("Can't get local account.")
		os.Exit(1)
	}
	acc, err := cli.GetDefaultAccount()
	checkErr(err,"can't get main-account")
	rpchttp.Wallet = cli
	//init node server for sync data
	node := net.StartProtocol(acc.PublicKey)
	rpchttp.RegistRpcNode(node)
	time.Sleep(6 * time.Second)
	log.Info("start sync block")
	node.SyncNodeHeight()
	log.Info("sync block finish")
	go node.WaitForFourPeersStart()
	go node.WaitForSyncBlkFinish()
	//start rpc server for console
	go rpchttp.StartRPCServer()
	//start server for http api
	go httprestful.StartServer(node)
	//if this is verity node ,start consensus protocol
	if protocol.VERIFYNODENAME == config.Parameters.NodeType {
		dbftServices := dbft.NewDbftService(cli, "logcon", node)
		rpchttp.RegistDbftService(dbftServices)
		go dbftServices.Start()
	}
	for {
		time.Sleep(dbft.GenBlockTime)
		log.Info("BlockHeight = ", ledger.DefaultLedger.Blockchain.BlockHeight)
		//isNeedNewFile := log.CheckIfNeedNewFile()
		//if isNeedNewFile == true {
		//	log.ClosePrintLog()
		//	log.Init(log.Path, os.Stdout)
		//}
	}
}


func checkErr(err error,msg string){
	if err != nil {
		if msg == ""{
			log.Error(err)
		}
		log.Error(msg)
		os.Exit(1)
	}
}