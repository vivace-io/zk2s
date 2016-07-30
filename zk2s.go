package main

import (
	"os"

	"github.com/eveopsec/zk2s/zk2s"
	"github.com/eveopsec/zk2s/zk2s/config"
	"github.com/urfave/cli"
)

const VERSION = "0.5"

var CONTRIBUTORS = []cli.Author{
	cli.Author{
		Name: "Nathan \"Vivace Naaris\" Morley",
	},
	cli.Author{
		Name: "\"Zuke\"",
	},
}

func main() {
	app := cli.NewApp()
	app.Authors = CONTRIBUTORS
	app.Version = VERSION
	app.Name = "zk2s"
	app.Usage = "A Slack bot for posting kills from zKillboard to slack in near-real time."
	app.Commands = []cli.Command{
		zk2s.CMD_Run,
		config.CMD_Config,
	}
	app.Run(os.Args)
}
