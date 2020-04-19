package bot

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
)

type Plugin struct {
	Logger *logrus.Logger
}

func (p Plugin) ReceiveMessage(input *botplug.MessageInput) (output *botplug.MessageOutput) {
	var queue []interface{}

	p.Logger.WithFields(logrus.Fields{
		"source": *input.Source,
	}).Debug(input.Messages)

	// TODO
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
