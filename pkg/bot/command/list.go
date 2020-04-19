package command

import (
	"fmt"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

const (
	command = `/list`
)

type PluginList struct {
	SharedMem *sharedmem.SharedMem
	Rcon      *rcon.Client
	Logger    *logrus.Logger
}

func (p PluginList) CommandName() string {
	return command
}

func (p PluginList) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	// read data from SharedMem
	data, err := p.SharedMem.ReadSharedMem()
	if err != nil {
		p.Logger.Warn(err)
		queue = append(queue, "Internal Error")
		return &botplug.MessageOutput{Queue: queue}
	}

	// assertion to domain.Domain
	domainData, ok := data.(domain.Domain)
	if !ok {
		p.Logger.Warn(fmt.Errorf("data cannot assertion to domain.Domain"))
	}

	// ログイン中のユーザを LINE に送信
	var loginUsernames []string
	for _, user := range domainData.LoginUsers {
		loginUsernames = append(loginUsernames, user.Name)
	}
	queue = append(queue, fmt.Sprintf("ログイン中のユーザ: %s", loginUsernames))
	return &botplug.MessageOutput{Queue: queue}
}
