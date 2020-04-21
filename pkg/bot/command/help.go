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

	msg := `
/help
ヘルプメッセージを表示します

/list
ログイン中のユーザ一覧を表示します

/title hoge
Minecraftのゲーム画面に hoge と表示されます

/whitelist list
ホワイトリストを表示します

/whitelist add hoge
ユーザ hoge をホワイトリストに追加します

/whitelist delete hoge
ユーザ hoge をホワイトリストから削除します
`
	queue = append(queue, msg)

	return &botplug.MessageOutput{Queue: queue}
}
