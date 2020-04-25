package localmem

import (
	"fmt"
	"sync"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
	"github.com/sirupsen/logrus"
)

var (
	mu           sync.Mutex
	sharedMemory *domain.Entity
)

type SharedMem struct {
	logger        *logrus.Logger
	sendStream    chan<- domain.Entity
	receiveStream <-chan domain.Entity
}

func New(logger *logrus.Logger) *SharedMem {
	stream := make(chan domain.Entity)
	m := &SharedMem{
		logger:        logger,
		sendStream:    stream,
		receiveStream: stream,
	}
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
