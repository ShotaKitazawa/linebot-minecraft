package sharedmem

import "github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"

type SharedMem interface {
	ReadSharedMem() (*domain.Domain, error)
	SendToChannel(data domain.Domain) error
}
