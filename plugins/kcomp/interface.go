package kcomp

import "context"

type Message struct {
	Key     []byte
	Payload []byte
}

type KProducer interface {
	Publish(ctx context.Context, topic string, key string, value interface{}) error
}

type KConsumer interface {
	Subscribe(groupId string, topic string, handlerFunc func(msg *Message) error)
}

type HandlerFunc func(msg *Message) error
