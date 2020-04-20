package command

import (
	"fmt"
	"strings"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type PluginTitle struct {
	SharedMem *sharedmem.SharedMem
	Rcon      *rcon.Client
	Logger    *logrus.Logger
}

func (p PluginTitle) CommandName() string {
	return `/title`
}

func (p PluginTitle) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	// send RCON
	destUsers, err := p.Rcon.Title(strings.Join(input.Messages[1:], " "))
	if err != nil {
		p.Logger.Error(err)
		queue = append(queue, "Internal Error")
		return &botplug.MessageOutput{Queue: queue}
	}
	if len(destUsers) == 0 {
		queue = append(queue, `ログイン中のユーザは存在しません`)
		return &botplug.MessageOutput{Queue: queue}
	}
	for _, user := range destUsers {
		queue = append(queue, fmt.Sprintf(`%s に送信しました`, user))
	}

	return &botplug.MessageOutput{Queue: queue}
}
