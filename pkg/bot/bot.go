package bot

import (
	"strconv"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	Logger    *logrus.Logger
	SharedMem *sharedmem.SharedMem
}

func (p *Plugin) ReceiveMessage(input *botplug.MessageInput) *botplug.MessageOutput {
	var queue []interface{}

	p.Logger.WithFields(logrus.Fields{
		"source": *input.Source,
	}).Debug(input.Messages)

	data, err := p.SharedMem.ReadSharedMem()
	if err != nil {
		queue = append(queue, "Internal Error")
		p.Logger.Warn(err)
		return &botplug.MessageOutput{Queue: queue}
	}

	queue = append(queue, strconv.Itoa(data.(int)))
	return &botplug.MessageOutput{Queue: queue}
}
