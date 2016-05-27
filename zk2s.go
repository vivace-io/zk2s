package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/eveopsec/zk2s/util"
	"github.com/nlopes/slack"
	"github.com/vivace-io/evelib/zkill"
	"github.com/vivace-io/gonfig"
)

/* zk2s.go
 * Main entrypoint and controller for zk2s
 */

const VERSION = "0.3"

var CONTRIBUTORS = []cli.Author{
	cli.Author{
		Name: "Nathan \"Vivace Naaris\" Morley",
	},
	cli.Author{
		Name: "\"Zuke\"",
	},
}

var config *util.Configuration
var bot *slack.Client

func main() {
	app := cli.NewApp()
	app.Authors = CONTRIBUTORS
	app.Version = VERSION
	app.Name = "zk2s"
	app.Usage = "a Slack bot for posting kills from zKillboard to slack in near-real time"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "start",
			Usage:  "start zk2s application",
			Action: Run,
		},
		cli.Command{
			Name:   "configure",
			Usage:  "configure zk2s application to be run",
			Action: util.RunConfigure,
		},
	}
	app.Run(os.Args)
}

// Run zk2s
func Run(c *cli.Context) {
	log.Printf("%v version %v", c.App.Name, c.App.Version)
	var err error

	// 1 - Load Configuration file
	config = new(util.Configuration)
	err = gonfig.Load(config)
	if err != nil {
		log.Fatalf("Unable to read config with error %v", err)
		os.Exit(1)
	}
	// 2 - Setup a new Slack Bot
	bot = slack.New(config.BotToken)
	authResp, err := bot.AuthTest()
	if err != nil {
		log.Fatalf("Unable to authenticate with Slack - %v", err)
		os.Exit(1)
	}
	log.Printf("Connected to Slack Team %v as user %v", authResp.Team, authResp.User)

	// 3 - Watch for new kills and log errors
	errc := make(chan error, 5)
	killc := make(chan zkill.Kill, 10)
	zClient := zkill.NewRedisQ()
	zClient.UserAgent = config.UserAgent
	zClient.FetchKillmails(killc, errc)
	handleKills(killc)
	handleErrors(errc)
	select {}
}

// handleKills sends the kill to be filtered/processed before posting to slack.
func handleKills(killChan chan zkill.Kill) {
	go func() {
		for {
			select {
			case kill := <-killChan:
				PostKill(&kill)
			}
		}
	}()
}

// handleErrors logs errors returned in Zkillboard queries
func handleErrors(errChan chan error) {
	go func() {
		for {
			select {
			case err := <-errChan:
				log.Printf("ERROR - %v", err.Error())
			}
		}
	}()
}
