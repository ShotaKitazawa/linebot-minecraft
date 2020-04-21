package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/bot"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug/line"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/eventer"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

var (
	// These variables are set in build step
	Version  = "unset"
	Revision = "unset"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
}

type argsConfig struct {
	channelSecret string
	channelToken  string
	groupID       string
	rconHost      string
	rconPort      int
	rconPassword  string
}

func newArgsConfig() *argsConfig {
	cfg := &argsConfig{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.channelSecret, "line-channel-secret", "", "")
	fl.StringVar(&cfg.channelToken, "line-channel-token", "", "")
	fl.StringVar(&cfg.groupID, "line-group-id", "", "specified LINE Group ID, send push message to this Group")
	fl.StringVar(&cfg.rconHost, "rcon-host", "", "RCON Host")
	fl.IntVar(&cfg.rconPort, "rcon-port", 25575, "RCON Port")
	fl.StringVar(&cfg.rconPassword, "rcon-password", "", "RCON Password")

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show application version")
	fl.Parse(os.Args[1:])

	if showVersion {
		fmt.Printf("version: %s (revision %s)", Version, Revision)
	}

	if cfg.channelSecret == "" ||
		cfg.channelToken == "" ||
		cfg.groupID == "" ||
		cfg.rconHost == "" ||
		cfg.rconPort == 0 ||
		cfg.rconPassword == "" {
		fmt.Println("not enough required fields")
		os.Exit(2)
	}

	return cfg
}

func main() {
	var err error

	// initialize
	args := newArgsConfig()
	m := sharedmem.New()
	rcon, err := rcon.New(args.rconHost, args.rconPort, args.rconPassword)
	if err != nil {
		panic(err)
	}
	eventer, err := eventer.New(args.groupID, args.channelSecret, args.channelToken, m, rcon, logger)
	if err != nil {
		panic(err)
	}

	// run eventer
	go eventer.Run()

	// TODO: run exporter

	// run bot
	handler, err := line.NewHandler(&line.Config{
		GroupID:       args.groupID,
		ChannelSecret: args.channelSecret,
		ChannelToken:  args.channelToken,
		Plugin:        bot.New(m, rcon, logger),
	})
	if err != nil {
		panic(err)
	}
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
