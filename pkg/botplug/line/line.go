package line

import (
	"log"
	"net/http"
	"strings"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

type Config struct {
	GroupID       string
	ChannelSecret string
	ChannelToken  string
	Plugin        botplug.BotPlugin
}

func NewHandler(config Config) (*httphandler.WebhookHandler, error) {
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

				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					input := &botplug.MessageInput{
						Timestamp: event.Timestamp,
						Source: &botplug.Source{
							Type:    string(event.Source.Type),
							UserID:  event.Source.UserID,
							GroupID: event.Source.GroupID,
						},
						Messages: strings.Split(message.Text, " "),
					}

					output := config.Plugin.RecieveMessage(input)

					for _, element := range output.Queue {
						switch typedElement := element.(type) {
						case string:
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(typedElement)).Do(); err != nil {
								log.Print(err)
							}
						case []linebot.SendingMessage:
							_, err = bot.PushMessage(config.GroupID, typedElement...).Do()
							if err != nil {
								log.Print(err)
							}
						}
					}
				}
			}
		}
	})
	return handler, nil
}
