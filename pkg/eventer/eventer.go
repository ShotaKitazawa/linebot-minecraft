package eventer

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

type Eventer struct {
	domain.LineClientConfig
	sharedMem *sharedmem.SharedMem
	rcon      *rcon.Client
	Logger    *logrus.Logger
}

func New(groupID, channelSecret, channelToken string, m *sharedmem.SharedMem, rcon *rcon.Client, logger *logrus.Logger) *Eventer {
	return &Eventer{
		LineClientConfig: domain.LineClientConfig{
			GroupID:       groupID,
			ChannelSecret: channelSecret,
			ChannelToken:  channelToken,
		},
		sharedMem: m,
		rcon:      rcon,
		Logger:    logger,
	}
}

func (e *Eventer) Run() error {
	return e.cronjob()
}

func (e *Eventer) cronjob() error {
	if err := e.job(); err != nil {
		e.Logger.Error(err)
	}
	t := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-t.C:
			if err := e.job(); err != nil {
				e.Logger.Error(err)
			}
		}
	}
	// t.Stop()
	// return nil
}

func (e *Eventer) job() error {
	var err error
	var d domain.Domain

	// get Minecraft metrics by RCON
	loginUsernames, err := e.rcon.List()
	if err != nil {
		return err
	}
	for _, username := range loginUsernames {
		userData, err := e.rcon.DataGetEntity(username)
		if err != nil {
			return err
		}
		d.LoginUsers = append(d.LoginUsers, domain.User{
			Name:    username,
			XpLevel: userData.XpLevel,
			Position: domain.Position{
				X: userData.X,
				Y: userData.Y,
				Z: userData.Z,
			},
		})
	}
	d.WhitelistUsernames, err = e.rcon.WhitelistList()
	if err != nil {
		return err
	}

	// TODO: send to LINE (PUSH notification) if d.LoginUsers != sharedmem.Domain.LoginUsers

	// write to chan
	e.sharedMem.SendToChannel(d)

	return nil
}
