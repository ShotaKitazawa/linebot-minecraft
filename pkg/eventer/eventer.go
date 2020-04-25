package eventer

import (
	"fmt"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

const (
	cronJobInterval = 10
)

type Eventer struct {
	domain.LineClientConfig

	sharedMem *sharedmem.SharedMem
	rcon      *rcon.Client
	Logger    *logrus.Logger
}

func New(groupIDs, channelSecret, channelToken string, m *sharedmem.SharedMem, rcon *rcon.Client, logger *logrus.Logger) (*Eventer, error) {
	client, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}
	return &Eventer{
		LineClientConfig: domain.LineClientConfig{
			GroupIDs: strings.Split(groupIDs, ","),
			Client:   client,
		},
		sharedMem: m,
		rcon:      rcon,
		Logger:    logger,
	}, nil
}

func (e *Eventer) Run() error {
	return e.cronjob()
}

func (e *Eventer) cronjob() error {
	if err := e.job(); err != nil {
		e.Logger.Error(err)
	}
	t := time.NewTicker(cronJobInterval * time.Second)
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
	currentLoginUserSet := mapset.NewSet()
	currentLoginUsernames, err := e.rcon.List()
	if err != nil {
		return err
	}
	for _, username := range currentLoginUsernames {
		userData, err := e.rcon.DataGetEntity(username)
		if err != nil {
			return err
		} else if userData == nil {
			e.Logger.Warn(`userData is nil`)
			return nil
		}
		currentLoginUser := domain.User{
			Name:    username,
			Health:  userData.Health,
			XpLevel: userData.XpLevel,
			Position: domain.Position{
				X: userData.X,
				Y: userData.Y,
				Z: userData.Z,
			},
		}
		d.LoginUsers = append(d.LoginUsers, currentLoginUser)
		currentLoginUserSet.Add(currentLoginUser.Name)
	}
	d.WhitelistUsernames, err = e.rcon.WhitelistList()
	if err != nil {
		return err
	}

	// get logged in users from SharedMem
	previousLoginUserSet := mapset.NewSet()
	data, err := e.sharedMem.ReadSharedMem()
	if err != nil {
		// write to sharedMem & return
		e.sharedMem.SendToChannel(d)
		return err
	}
	for _, previousLoginUser := range data.LoginUsers {
		previousLoginUserSet.Add(previousLoginUser.Name)
	}

	// send to LINE (PUSH notification) if d.LoginUsers != sharedmem.Domain.LoginUsers
	loggingInUsernameSet := currentLoginUserSet.Difference(previousLoginUserSet)
	if loggingInUsernameSet.Cardinality() != 0 {
		for _, groupID := range e.GroupIDs {
			if _, err := e.Client.PushMessage(groupID, linebot.NewTextMessage(fmt.Sprintf(`ユーザがログインしました: %v`, loggingInUsernameSet))).Do(); err != nil {
				e.Logger.Error(`failed to push notification: `, err)
			}
		}
	}
	loggingOutUsernameSet := previousLoginUserSet.Difference(currentLoginUserSet)
	if loggingOutUsernameSet.Cardinality() != 0 {
		for _, groupID := range e.GroupIDs {
			if _, err := e.Client.PushMessage(groupID, linebot.NewTextMessage(fmt.Sprintf(`ユーザがログアウトしました: %v`, loggingOutUsernameSet))).Do(); err != nil {
				e.Logger.Error(`failed to push notification: `, err)
			}
		}
	}

	// write to sharedMem
	e.sharedMem.SendToChannel(d)

	return nil
}
