package bot

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/bot/command"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

type PluginConfig struct {
	SharedMem *sharedmem.SharedMem
	Rcon      *rcon.Client
	Logger    *logrus.Logger
	Plugins   []PluginInterface
}

func New(m *sharedmem.SharedMem, rcon *rcon.Client, logger *logrus.Logger) *PluginConfig {
	return &PluginConfig{
		SharedMem: m,
		Rcon:      rcon,
		Logger:    logger,
		Plugins: []PluginInterface{
			command.PluginList{
				SharedMem: m,
				Rcon:      rcon,
				Logger:    logger,
			},
			command.PluginTitle{
				SharedMem: m,
				Rcon:      rcon,
				Logger:    logger,
			},
			command.PluginWhitelist{
				SharedMem: m,
				Rcon:      rcon,
				Logger:    logger,
			},
			command.PluginHelp{
				SharedMem: m,
				Rcon:      rcon,
				Logger:    logger,
			},
		},
	}
}

func (pc *PluginConfig) ReceiveMessageEntry(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	if !strings.HasPrefix(input.Messages[0], "/") {
		return nil
	}

	pc.Logger.WithFields(logrus.Fields{
		"source": *input.Source,
	}).Debug(input.Messages)

	for _, plugin := range pc.Plugins {
		if input.Messages[0] == plugin.CommandName() {
			return plugin.ReceiveMessage(input)
		}
	}

	queue = append(queue, `no such command`)
	return &botplug.MessageOutput{Queue: queue}
}
