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
	// TODO typed: interface{} -> []domain.XXX{}
	data interface{}
}

func New() *SharedMem {
	stream := make(chan domain.Domain)
	m := new(SharedMem)
	m.sendStream = stream
	m.receiveStream = stream
	go m.receiveFromChannelAndWriteSharedMem()
	return m
}

// TODO typed: interface{} -> domain.XXX{}
func (m *SharedMem) ReadSharedMem() (interface{}, error) {
	var result interface{}
	mu.Lock()
	result = m.data
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
			m.data = d
			mu.Unlock()
		}
	}
	// return nil
}
