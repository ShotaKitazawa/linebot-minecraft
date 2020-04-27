package eventer

import (
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/botplug"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain/i18n"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/rcon"
	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

const (
	cronJobInterval = 10
)

type Eventer struct {
	botplug.BotPluginSender

	MinecraftHostname string
	sharedMem         sharedmem.SharedMem
	rcon              rcon.RconClient
	Logger            *logrus.Logger
}

func New(minecraftHostname string, sender botplug.BotPluginSender, m sharedmem.SharedMem, rcon rcon.RconClient, logger *logrus.Logger) (*Eventer, error) {
	return &Eventer{
		BotPluginSender:   sender,
		MinecraftHostname: minecraftHostname,
		sharedMem:         m,
		rcon:              rcon,
		Logger:            logger,
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
	var currentData domain.Entity

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
		currentData.LoginUsers = append(currentData.LoginUsers, currentLoginUser)
		currentLoginUserSet.Add(currentLoginUser.Name)
	}
	currentData.WhitelistUsernames, err = e.rcon.WhitelistList()
	if err != nil {
		return err
	}

	// get logged in users from SharedMem
	previousLoginUserSet := mapset.NewSet()
	previousData, err := e.sharedMem.SyncReadEntityFromSharedMem()
	if err != nil {
		// write to sharedMem & return
		return e.sharedMem.AsyncWriteEntityToSharedMem(currentData)
	}
	for _, previousLoginUser := range previousData.LoginUsers {
		previousLoginUserSet.Add(previousLoginUser.Name)

		/*
			// check if user log out
			userIsLogOut := true
			for _, currentLoginUser := range currentData.LoginUsers {
				if previousLoginUser.Name == currentLoginUser.Name{
					userIsLogOut = false
				}
			}
			if userIsLogOut{
				currentData.AllUsers = append(currentData.AllUsers, previousLoginUser)
			}
		*/
	}

	// store domain.AllUsers, LogoutUsers
	for _, currentUser := range currentData.LoginUsers {
		currentData.AllUsers = append(currentData.AllUsers, currentUser)
	}
	for _, previousUser := range previousData.AllUsers {
		var flag bool
		for _, currentUser := range currentData.LoginUsers {
			if previousUser.Name == currentUser.Name {
				flag = true
			}
		}
		if !flag {
			currentData.AllUsers = append(currentData.AllUsers, previousUser)
		}
	}

	// send to LINE (PUSH notification) if d.LoginUsers != sharedmem.Domain.LoginUsers
	loggingInUsernameSet := currentLoginUserSet.Difference(previousLoginUserSet)
	if loggingInUsernameSet.Cardinality() != 0 {
		e.BotPluginSender.SendTextMessage(i18n.T.Sprintf(i18n.MessageUsersLogin, loggingInUsernameSet.ToSlice()))
	}
	loggingOutUsernameSet := previousLoginUserSet.Difference(currentLoginUserSet)
	if loggingOutUsernameSet.Cardinality() != 0 {
		e.BotPluginSender.SendTextMessage(i18n.T.Sprintf(i18n.MessageUsersLogout, loggingOutUsernameSet.ToSlice()))
	}

	// write to sharedMem
	e.sharedMem.AsyncWriteEntityToSharedMem(currentData)

	return nil
}
