package eventer

import (
	"time"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/sharedmem"
)

type Eventer struct {
	// TODO typed: interface{} -> domain.XXX{}
	sharedMem *sharedmem.SharedMem
}

// TODO: give config to connect to Redis
func New(m *sharedmem.SharedMem) *Eventer {
	return &Eventer{
		sharedMem: m,
	}
}

func (e *Eventer) Run() error {
	return e.stream()
}

func (e *Eventer) stream() error {
	var cnt int
	for {
		cnt++
		time.Sleep(time.Second)
		e.sharedMem.SendToChannel(cnt)
		// TODO
		// subscribe from Redis
		// send to LINE (push) if XXX
		// write to chan
	}
	//return nil
}
