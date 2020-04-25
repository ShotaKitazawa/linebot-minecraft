package sharedmem

import "github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"

type SharedMem interface {
	SyncReadEntityFromSharedMem() (*domain.Entity, error)
	AsyncWriteEntityToSharedMem(data domain.Entity) error
}
