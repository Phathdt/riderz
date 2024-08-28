package kcomp

type KProducer interface {
	Publish(topic string, key string, value interface{}) error
}
