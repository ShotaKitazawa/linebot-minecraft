package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug/line"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Plugin struct{}

func (p Plugin) RecieveMessage(input *botplug.MessageInput) (output *botplug.MessageOutput) {
	var queue []interface{}

	// test
	queue = append(queue, "test")
	leftBtn := linebot.NewMessageAction("left", "left clicked")
	rightBtn := linebot.NewMessageAction("right", "right clicked")
	template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
	message := linebot.NewTemplateMessage("Sorry :(, please update your app.", template)
	var messages []linebot.SendingMessage
	messages = append(messages, message)
	queue = append(queue, messages)

	return &botplug.MessageOutput{Queue: queue}
}

func main() {
	plugin := Plugin{}
	conf := line.Config{
		GroupID:       os.Getenv("GROUP_ID"),
		ChannelSecret: os.Getenv("CHANNEL_SECRET"),
		ChannelToken:  os.Getenv("CHANNEL_TOKEN"),
		Plugin:        plugin,
	}

	// /callback にエンドポイントの定義
	handler, err := line.NewHandler(conf)
	if err != nil {
		panic(err)
	}
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
