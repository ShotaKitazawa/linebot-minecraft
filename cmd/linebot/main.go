package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/bot"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug/line"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
}

func main() {
	conf := line.Config{
		GroupID:       os.Getenv("GROUP_ID"),
		ChannelSecret: os.Getenv("CHANNEL_SECRET"),
		ChannelToken:  os.Getenv("CHANNEL_TOKEN"),
		Plugin:        bot.Plugin{Logger: logger},
	}

	handler, err := line.NewHandler(conf)
	if err != nil {
		panic(err)
	}
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
