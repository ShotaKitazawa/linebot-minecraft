package sharedmem

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

type SharedMem struct {
	// TODO : sendChannel -> sendChannels
	sendStream chan<- interface{}
	// TODO : receiveChannel -> receiveChannels
	receiveStream <-chan interface{}
	// TODO typed: interface{} -> []domain.XXX{}
	data interface{}
}

func New() *SharedMem {
	stream := make(chan interface{})
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

func (m *SharedMem) SendToChannel(data interface{}) error {
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
