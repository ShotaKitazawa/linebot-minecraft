package bot

import "github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"

type PluginInterface interface {
	CommandName() string
	ReceiveMessage(*botplug.MessageInput) *botplug.MessageOutput
}
