package localmem

import (
	"fmt"
	"sync"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
)

var (
	mu           sync.Mutex
	sharedMemory *domain.Entity
)

type SharedMem struct {
	// TODO : sendChannel -> sendChannels
	sendStream chan<- domain.Entity
	// TODO : receiveChannel -> receiveChannels
	receiveStream <-chan domain.Entity
}

func New() *SharedMem {
	stream := make(chan domain.Entity)
	m := new(SharedMem)
	m.sendStream = stream
	m.receiveStream = stream
	go m.receiveFromChannelAndWriteSharedMem()
	return m
}

func (m *SharedMem) SyncReadEntityFromSharedMem() (*domain.Entity, error) {
	mu.Lock()
	result := sharedMemory
	mu.Unlock()
	if result == nil {
		return nil, fmt.Errorf("no such data")
	}
	return result, nil
}

func (m *SharedMem) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	m.sendStream <- data
	return nil
}

func (m *SharedMem) receiveFromChannelAndWriteSharedMem() error {
	for {
		select {
		case d := <-m.receiveStream:
			mu.Lock()
			sharedMemory = &d
			mu.Unlock()
		}
	}
	// return nil
}
