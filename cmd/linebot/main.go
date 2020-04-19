package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/bot"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug/line"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/eventer"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
}

func main() {
	var err error

	m := sharedmem.New()
	plugin := &bot.Plugin{Logger: logger, SharedMem: m}
	conf := &line.Config{
		GroupID:       os.Getenv("GROUP_ID"),
		ChannelSecret: os.Getenv("CHANNEL_SECRET"),
		ChannelToken:  os.Getenv("CHANNEL_TOKEN"),
		Plugin:        plugin,
	}
	eventer := eventer.New(m)
	go eventer.Run()

	handler, err := line.NewHandler(conf)
	if err != nil {
		panic(err)
	}
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
