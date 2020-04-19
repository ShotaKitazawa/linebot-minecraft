package sharedmem

import (
	"fmt"
	"sync"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
)

var mu sync.Mutex

type SharedMem struct {
	// TODO : sendChannel -> sendChannels
	sendStream chan<- domain.Domain
	// TODO : receiveChannel -> receiveChannels
	receiveStream <-chan domain.Domain
	data          *domain.Domain
}

func New() *SharedMem {
	stream := make(chan domain.Domain)
	m := new(SharedMem)
	m.sendStream = stream
	m.receiveStream = stream
	go m.receiveFromChannelAndWriteSharedMem()
	return m
}

func (m *SharedMem) ReadSharedMem() (*domain.Domain, error) {
	mu.Lock()
	result := m.data
	mu.Unlock()
	if result == nil {
		return nil, fmt.Errorf("no such data")
	}
	return result, nil
}

func (m *SharedMem) SendToChannel(data domain.Domain) error {
	m.sendStream <- data
	return nil
}

func (m *SharedMem) receiveFromChannelAndWriteSharedMem() error {
	for {
		select {
		case d := <-m.receiveStream:
			mu.Lock()
			m.data = &d
			mu.Unlock()
		}
	}
	// return nil
}
