package kcomp

import (
	"context"
	"time"
)

type HandlerFunc func(msg *Message) error
type HandlerMultiFunc func(msgs []*Message) error

type Message struct {
	Key     []byte
	Payload []byte
}

type KProducer interface {
	Publish(ctx context.Context, topic string, key string, value interface{}) error
}

type KConsumer interface {
	Subscribe(groupId string, topic string, handlerFunc HandlerFunc)
	BatchSubscribe(groupId string, topic string, batchTimeout time.Duration, batchSize int, handlerFunc HandlerMultiFunc)
}
