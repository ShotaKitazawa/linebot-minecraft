package redis

import (
	"encoding/json"
	"strconv"

	"github.com/ShotaKitazawa/linebot-minecraft/pkg/domain"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type SharedMem struct {
	logger        *logrus.Logger
	sendStream    chan<- domain.Entity
	receiveStream <-chan domain.Entity
	Conn          redis.Conn
	redisHostname string
}

func New(logger *logrus.Logger, addr string, port int) (*SharedMem, error) {
	stream := make(chan domain.Entity)
	redisHostname := addr + ":" + strconv.Itoa(port)
	c, err := redis.Dial("tcp", redisHostname)
	if err != nil {
		return nil, err
	}
	m := &SharedMem{
		logger:        logger,
		sendStream:    stream,
		receiveStream: stream,
		Conn:          c,
		redisHostname: redisHostname,
	}
	go m.receiveFromChannelAndWriteSharedMem()
	return m, nil
}

func (m *SharedMem) SyncReadEntityFromSharedMem() (*domain.Entity, error) {
	data, err := redis.Bytes(m.Conn.Do("GET", "entity"))
	if err != nil {
		m.logger.Warn(err)
		return nil, m.reconnect()
	} else if data == nil {
		return nil, nil
	}
	entity := domain.Entity{}
	if err := json.Unmarshal(data, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (m *SharedMem) AsyncWriteEntityToSharedMem(data domain.Entity) error {
	m.sendStream <- data
	return nil
}

func (m *SharedMem) receiveFromChannelAndWriteSharedMem() error {
	for {
		select {
		case d := <-m.receiveStream:
			data, err := json.Marshal(&d)
			if err != nil {
				return err
			}
			_, err = m.Conn.Do("SET", "entity", data)
			if err != nil {
				m.logger.Warn(err)
				return m.reconnect()
			}
		}
	}
	// return nil
}

func (m *SharedMem) reconnect() error {
	c, err := redis.Dial("tcp", m.redisHostname)
	if err != nil {
		return err
	}
	m.Conn = c
	return nil
}
