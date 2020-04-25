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
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem/localmem"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem/redis"
)

var (
	// These variables are set in build step
	Version  = "unset"
	Revision = "unset"
)

var logger = logrus.New()

func newLogger(loglevel string) *logrus.Logger {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	switch loglevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		panic(fmt.Errorf("newLogger: invalid arguments"))
	}
	return logger
}

type argsConfig struct {
	loglevel      string
	channelSecret string
	channelToken  string
	groupIDs      string
	sharedmemMode string
	redisHost     string
	redisPort     int
	rconHost      string
	rconPort      int
	rconPassword  string
}

func newArgsConfig() *argsConfig {
	cfg := &argsConfig{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.loglevel, "log-level", "info", "Log Level (debug, info, warn, error)")
	fl.StringVar(&cfg.channelSecret, "line-channel-secret", "", "LINE Bot's Channel Secret")
	fl.StringVar(&cfg.channelToken, "line-channel-token", "", "LINE Bot's Channel Token")
	fl.StringVar(&cfg.groupIDs, "line-group-id", "", "specified LINE Group ID, send push message to this Group")
	fl.StringVar(&cfg.sharedmemMode, "sharedmem-mode", "local", `using Shared Memory ("local" or "redis")`)
	fl.StringVar(&cfg.redisHost, "redis-host", "127.0.0.1", "Redis Host (enabled when sharedmem-mode=redis)")
	fl.IntVar(&cfg.redisPort, "redis-port", 6379, "Redis Port (enabled when sharedmem-mode=redis)")
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

	if !(cfg.sharedmemMode == "local" ||
		cfg.sharedmemMode == "redis") {
		fmt.Println("sharedmemMode mismatch")
		os.Exit(2)

	}

	if !(cfg.loglevel == "debug" ||
		cfg.loglevel == "info" ||
		cfg.loglevel == "warn" ||
		cfg.loglevel == "error") {
		fmt.Println("log-level mismatch")
		os.Exit(2)
	}

	return cfg
}

func main() {
	var err error

	// args parse
	args := newArgsConfig()

	// set logger
	logger = newLogger(args.loglevel)

	m := func(sharedmemMode string) sharedmem.SharedMem {
		switch sharedmemMode {
		case "local":
			// run sharedMem & get sharedMem instance
			m, err := localmem.New(logger)
			if err != nil {
				panic(err)
			}
			return m
		case "redis":
			m, err := redis.New(logger, args.redisHost, args.redisPort)
			if err != nil {
				panic(err)
			}
			return m
		default:
			panic(fmt.Errorf("sharedmemMode mismatch"))
		}
	}(args.sharedmemMode)

	// get rcon instance
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
