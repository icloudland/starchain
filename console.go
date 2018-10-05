package main

import (
	"github.com/urfave/cli"
	"github.com/icloudland/starchain/client"
	"github.com/icloudland/starchain/client/asset"
	"os"
	"github.com/icloudland/starchain/client/wallet"
	"github.com/icloudland/starchain/client/bookeeper"
	"github.com/icloudland/starchain/client/consensus"
	"github.com/icloudland/starchain/client/info"
	"github.com/icloudland/starchain/client/smartcontract"
)

func main(){
	app := cli.NewApp()
	app.Name = "stc-client"
	app.Version = "v0.0.1"
	app.HelpName = "stc-client"
	app.Usage = "command line tool for STC blockchain"
	app.UsageText = "stc-client [global options] command [command options] [args]"
	app.HideHelp = false
	app.HideVersion = false

	app.Flags = []cli.Flag{
		client.NewIpFlag(),
		client.NewPortFlag(),
	}
	app.Commands = []cli.Command{
		*asset.NewCommand(),
		*wallet.NewCommand(),
		*bookeeper.NewCommond(),
		*consensus.NewCommond(),
		*info.NewCommand(),
		*smartcontract.NewCommand(),
	}
	app.Run(os.Args)
}
