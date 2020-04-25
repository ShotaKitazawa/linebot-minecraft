package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/bot"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug/line"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/eventer"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/exporter"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem/localmem"
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
	groupIDs      string
	rconHost      string
	rconPort      int
	rconPassword  string
}

func newArgsConfig() *argsConfig {
	cfg := &argsConfig{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.channelSecret, "line-channel-secret", "", "")
	fl.StringVar(&cfg.channelToken, "line-channel-token", "", "")
	fl.StringVar(&cfg.groupIDs, "line-group-id", "", "specified LINE Group ID, send push message to this Group")
	fl.StringVar(&cfg.rconHost, "rcon-host", "", "RCON Host")
	fl.IntVar(&cfg.rconPort, "rcon-port", 25575, "RCON Port")
	fl.StringVar(&cfg.rconPassword, "rcon-password", "", "RCON Password")

	var showVersion bool
	fl.BoolVar(&showVersion, "v", false, "show application version")
	fl.Parse(os.Args[1:])

	if showVersion {
		fmt.Printf("version: %s (revision %s)", Version, Revision)
	}

	if cfg.channelSecret == "" ||
		cfg.channelToken == "" ||
		cfg.groupIDs == "" ||
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

	// args parse
	args := newArgsConfig()

	// run sharedMem
	m := localmem.New()
	rcon, err := rcon.New(args.rconHost, args.rconPort, args.rconPassword)
	if err != nil {
		panic(err)
	}

	// run eventer
	eventer, err := eventer.New(args.groupIDs, args.channelSecret, args.channelToken, m, rcon, logger)
	if err != nil {
		panic(err)
	}
	go eventer.Run()

	// run exporter
	collector, err := exporter.New(m, logger)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())

	// run bot
	handler, err := line.NewHandler(&line.Config{
		GroupIDs:      strings.Split(args.groupIDs, ","),
		ChannelSecret: args.channelSecret,
		ChannelToken:  args.channelToken,
		Plugin:        bot.New(m, rcon, logger),
	})
	if err != nil {
		panic(err)
	}
	http.Handle("/linebot", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
