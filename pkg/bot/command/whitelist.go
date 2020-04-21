package command

import (
	"fmt"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type PluginWhitelist struct {
	SharedMem *sharedmem.SharedMem
	Rcon      *rcon.Client
	Logger    *logrus.Logger
}

func (p PluginWhitelist) CommandName() string {
	return `/whitelist`
}

func (p PluginWhitelist) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	switch input.Messages[1] {
	case `add`:
		queue, _ = p.add(input.Messages[2:])
	case `delete`:
		queue, _ = p.delete(input.Messages[2:])
	case `list`:
		queue, _ = p.list()
	default:
		queue = append(queue, `no such command`)
	}

	return &botplug.MessageOutput{Queue: queue}
}

func (p PluginWhitelist) add(users []string) ([]interface{}, error) {
	var queue []interface{}
	for _, username := range users {
		if p.Rcon.WhitelistAdd(username) != nil {
			queue = append(queue, fmt.Sprintf(`ユーザ指定が間違っています: %s`, username))
		} else {
			queue = append(queue, fmt.Sprintf(`ユーザをホワイトリストに追加しました: %s`, username))
		}
	}
	return queue, nil
}

func (p PluginWhitelist) delete(users []string) ([]interface{}, error) {
	var queue []interface{}
	for _, username := range users {
		if p.Rcon.WhitelistRemove(username) != nil {
			queue = append(queue, fmt.Sprintf(`ユーザ指定が間違っています: %s`, username))
		} else {
			queue = append(queue, fmt.Sprintf(`ユーザをホワイトリストから削除しました: %s`, username))
		}
	}
	return queue, nil
}

func (p PluginWhitelist) list() ([]interface{}, error) {
	var queue []interface{}
	users, err := p.Rcon.WhitelistList()
	if err != nil {
		queue = append(queue, `Internal Error`)
		return nil, err
	}
	queue = append(queue, users)
	return queue, nil
}
