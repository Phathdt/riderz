package kconsumercomp

import (
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	sctx "github.com/phathdt/service-context"
	"riderz/plugins/kcomp"
	"sync"
)

type kConsumerComp struct {
	id       string
	brokers  string
	groupID  string
	logger   sctx.Logger
	consumer *kafka.Consumer
	handlers map[string]func(msg *kcomp.Message) error
	mu       sync.Mutex
	stopChan chan struct{}
}

func New(id string) *kConsumerComp {
	return &kConsumerComp{
		id:       id,
		handlers: make(map[string]func(msg *kcomp.Message) error),
		stopChan: make(chan struct{}),
	}
}

func (k *kConsumerComp) ID() string {
	return k.id
}

func (k *kConsumerComp) InitFlags() {
	flag.StringVar(&k.brokers, "kafka-brokers", "localhost:9092", "Kafka broker addresses, comma-separated")
	flag.StringVar(&k.groupID, "kafka-group-id", "my-group", "Kafka consumer group ID")
}

func (k *kConsumerComp) Activate(_ sctx.ServiceContext) error {
	k.logger = sctx.GlobalLogger().GetLogger(k.id)

	k.logger.Info("Connecting to Kafka...")

	config := &kafka.ConfigMap{
		"bootstrap.servers":  k.brokers,
		"group.id":           k.groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		k.logger.Error("Failed to create consumer", err.Error())
		return err
	}

	k.consumer = consumer

	go k.consumeMessages()

	k.logger.Info("Connected to Kafka")

	return nil
}

func (k *kConsumerComp) Stop() error {
	k.logger.Info("Closing Kafka consumer...")
	close(k.stopChan)
	return k.consumer.Close()
}

func (k *kConsumerComp) Subscribe(topic string, handlerFunc func(msg *kcomp.Message) error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.handlers[topic] = handlerFunc

	err := k.consumer.Subscribe(topic, nil)
	if err != nil {
		k.logger.Error("Failed to subscribe to topic", "topic", topic, "error", err.Error())
	}
}

func (k *kConsumerComp) consumeMessages() {
	for {
		select {
		case <-k.stopChan:
			return
		default:
			ev := k.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				msg := &kcomp.Message{
					Key:     e.Key,
					Payload: e.Value,
				}

				handler, ok := k.handlers[*e.TopicPartition.Topic]
				if !ok {
					k.logger.Error("Handler not found for topic", "topic", *e.TopicPartition.Topic)
					continue
				}

				func() {
					defer func() {
						if r := recover(); r != nil {
							k.logger.Error("Panic occurred while handling message", "recover", r)
						}
					}()

					err := handler(msg)
					if err != nil {
						k.logger.Error("Error handling message", "error", err.Error())
					}
				}()

				_, err := k.consumer.CommitMessage(e)
				if err != nil {
					k.logger.Error("Failed to commit message", "error", err.Error())
				}

			case kafka.Error:
				k.logger.Error("Kafka error", "error", e.String())

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}
