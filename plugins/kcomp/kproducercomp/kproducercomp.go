package kproducercomp

import (
	"encoding/json"
	"flag"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	sctx "github.com/phathdt/service-context"
)

type kProducerComp struct {
	id       string
	brokers  string
	logger   sctx.Logger
	producer *kafka.Producer
}

func New(id string) *kProducerComp {
	return &kProducerComp{id: id}
}

func (k *kProducerComp) ID() string {
	return k.id
}

func (k *kProducerComp) InitFlags() {
	flag.StringVar(
		&k.brokers,
		"kafka-brokers",
		"localhost:9092",
		"Kafka broker addresses, comma-separated",
	)
}

func (k *kProducerComp) Activate(_ sctx.ServiceContext) error {
	k.logger = sctx.GlobalLogger().GetLogger(k.id)

	k.logger.Info("Connecting to Kafka...")

	config := &kafka.ConfigMap{
		"bootstrap.servers": k.brokers,
		"client.id":         k.id,
		"acks":              "all",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		k.logger.Error("Failed to create producer", err.Error())
		return err
	}

	// Start a goroutine to handle delivery reports
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					k.logger.Error("Failed to deliver message", ev.TopicPartition.Error.Error())
				} else {
					k.logger.Info("Message delivered",
						"topic", *ev.TopicPartition.Topic,
						"partition", ev.TopicPartition.Partition,
						"offset", ev.TopicPartition.Offset)
				}
			}
		}
	}()

	k.producer = producer

	k.logger.Info("Connected to Kafka")

	return nil
}

func (k *kProducerComp) Stop() error {
	k.logger.Info("Closing Kafka producer...")
	k.producer.Flush(15 * 1000) // Wait for messages to be delivered
	k.producer.Close()
	return nil
}

func (k *kProducerComp) PublishMessage(topic string, key []byte, value []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	if key != nil && len(key) > 0 {
		message.Key = key
	}

	err := k.producer.Produce(message, nil)
	if err != nil {
		k.logger.Error("Failed to produce message", err.Error())
		return err
	}

	return nil
}

func (k *kProducerComp) Publish(topic string, key string, value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		k.logger.Error("Failed to serialize message", err.Error())
		return err
	}

	return k.PublishMessage(topic, []byte(key), payload)
}
