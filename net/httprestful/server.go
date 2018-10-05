package httprestful

import (
	"github.com/icloudland/starchain/net/protocol"
	"github.com/icloudland/starchain/net/httprestful/common"
	"github.com/icloudland/starchain/core/ledger"
	"strconv"
	."github.com/icloudland/starchain/common/config"
	"github.com/icloudland/starchain/events"
	"github.com/icloudland/starchain/net/httprestful/restful"
)

func StartServer(n protocol.Noder){
	common.SetNode(n)
	//events.NewEvent().AddListener(events.EventBlockPersistCompleted, SendBlock2NoticeServer)
	ledger.DefaultLedger.Blockchain.BCEvents.Subscribe(events.EventBlockPersistCompleted, SendBlock2NoticeServer)
	func() {
		rest := restful.InitRestServer(common.CheckAccessToken)
		go rest.Start()
	}()
}


func SendBlock2NoticeServer(v interface{}) {

	if len(Parameters.NoticeServerUrl) == 0 || !common.CheckPushBlock() {
		return
	}
	go func() {
		req := make(map[string]interface{})
		req["Height"] = strconv.FormatInt(int64(ledger.DefaultLedger.Blockchain.BlockHeight), 10)
		req = common.GetBlockByHeight(req)

		repMsg, _ := common.PostRequest(req, Parameters.NoticeServerUrl)
		if repMsg[""] == nil {
			//TODO
		}
	}()
}
