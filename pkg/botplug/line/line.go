package line

import (
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
)

type Config struct {
	GroupID       string
	ChannelSecret string
	ChannelToken  string
	Plugin        botplug.BotPlugin
}

func NewHandler(config *Config) (*httphandler.WebhookHandler, error) {
	handler, err := httphandler.New(
		config.ChannelSecret,
		config.ChannelToken,
	)
	if err != nil {
		return nil, err
	}
	bot, err := linebot.New(
		config.ChannelSecret,
		config.ChannelToken,
	)
	if err != nil {
		return nil, err
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch event.Message.(type) {
				case *linebot.TextMessage:
					if err = ReceiveTextMessage(event, bot, config); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	return handler, nil
}

func ReceiveTextMessage(event *linebot.Event, bot *linebot.Client, config *Config) (err error) {
	message := event.Message.(*linebot.TextMessage)
	input := &botplug.MessageInput{
		Timestamp: event.Timestamp,
		Source: &botplug.Source{
			Type:    string(event.Source.Type),
			UserID:  event.Source.UserID,
			GroupID: event.Source.GroupID,
		},
		Messages: strings.Split(message.Text, " "),
	}

	// execute user function
	output := config.Plugin.ReceiveMessageEntry(input)

	if output == nil {
		return
	}

	// proceed contents in queue
	for _, element := range output.Queue {
		switch typedElement := element.(type) {
		case string:
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(typedElement)).Do(); err != nil {
				return
			}
		case []linebot.SendingMessage:
			_, err = bot.PushMessage(config.GroupID, typedElement...).Do()
			if err != nil {
				return
			}
		}
	}
	return nil
}
