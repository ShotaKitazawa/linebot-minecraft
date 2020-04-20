package command

import (
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type PluginHelp struct {
	SharedMem *sharedmem.SharedMem
	Rcon      *rcon.Client
	Logger    *logrus.Logger
}

func (p PluginHelp) CommandName() string {
	return `/help`
}

func (p PluginHelp) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	msg := `ヘルプメッセージを表示します

/list
ログイン中のユーザ一覧を取得

/title hoge
Minecraftのゲーム画面に hoge と表示されます
`
	queue = append(queue, msg)

	return &botplug.MessageOutput{Queue: queue}
}
