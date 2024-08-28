package kcomp

type Message struct {
	Key     []byte
	Payload []byte
}

type KProducer interface {
	Publish(topic string, key string, value interface{}) error
}

type KConsumer interface {
	Subscribe(topic string, handlerFunc func(msg *Message) error)
}
