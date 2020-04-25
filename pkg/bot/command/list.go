package command

import (
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type PluginList struct {
	SharedMem sharedmem.SharedMem
	Logger    *logrus.Logger
}

func (p PluginList) CommandName() string {
	return `/list`
}

func (p PluginList) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	// read data from SharedMem
	data, err := p.SharedMem.SyncReadEntityFromSharedMem()
	if err != nil {
		p.Logger.Error(err)
		queue = append(queue, "Internal Error")
		return &botplug.MessageOutput{Queue: queue}
	}

	// ログイン中のユーザを LINE に送信
	var loginUsernames []string
	for _, user := range data.LoginUsers {
		loginUsernames = append(loginUsernames, user.Name)
	}
	if loginUsernames == nil {
		queue = append(queue, `ユーザが存在しません`)
		return &botplug.MessageOutput{Queue: queue}
	}
	queue = append(queue, loginUsernames)
	return &botplug.MessageOutput{Queue: queue}
}
