package command

import (
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain/i18n"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type PluginWhitelist struct {
	SharedMem sharedmem.SharedMem
	Rcon      *rcon.Client
	Logger    *logrus.Logger
}

func (p PluginWhitelist) CommandName() string {
	return `/whitelist`
}

func (p PluginWhitelist) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	if len(input.Messages) == 1 {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageInvalidArguments))
		return &botplug.MessageOutput{Queue: queue}
	}
	switch input.Messages[1] {
	case `add`:
		queue, _ = p.add(input.Messages[2:])
	case `delete`:
		queue, _ = p.delete(input.Messages[2:])
	case `list`:
		queue, _ = p.list()
	default:
		queue = append(queue, i18n.T.Sprintf(i18n.MessageInvalidArguments))
	}

	return &botplug.MessageOutput{Queue: queue}
}

func (p PluginWhitelist) add(users []string) ([]interface{}, error) {
	var queue []interface{}
	for _, username := range users {
		if p.Rcon.WhitelistAdd(username) != nil {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageUserIncorrect))
		} else {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageWhitelistAdd, username))
		}
	}
	return queue, nil
}

func (p PluginWhitelist) delete(users []string) ([]interface{}, error) {
	var queue []interface{}
	for _, username := range users {
		if p.Rcon.WhitelistRemove(username) != nil {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageUserIncorrect))
		} else {
			queue = append(queue, i18n.T.Sprintf(i18n.MessageWhitelistRemove, username))
		}
	}
	return queue, nil
}

func (p PluginWhitelist) list() ([]interface{}, error) {
	var queue []interface{}

	// read data from SharedMem
	data, err := p.SharedMem.SyncReadEntityFromSharedMem()
	if err != nil {
		p.Logger.Error(err)
		queue = append(queue, i18n.T.Sprintf(i18n.MessageError))
		return nil, err
	}

	// whitelist にいるユーザを LINE に送信
	var usernames []string
	for _, username := range data.WhitelistUsernames {
		usernames = append(usernames, username)
	}
	if usernames == nil {
		queue = append(queue, i18n.T.Sprintf(i18n.MessageNoUserExists))
		return queue, nil
	}
	queue = append(queue, usernames)
	return queue, nil
}
