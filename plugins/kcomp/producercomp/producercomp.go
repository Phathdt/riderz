package producercomp

import (
	"context"
	"encoding/json"
	"flag"
	"strings"
	"time"

	sctx "github.com/phathdt/service-context"
	"github.com/segmentio/kafka-go"
)

type producerComp struct {
	id      string
	brokers string
	logger  sctx.Logger
	writer  *kafka.Writer
}

func New(id string) *producerComp {
	return &producerComp{id: id}
}

func (c *producerComp) ID() string {
	return c.id
}

func (c *producerComp) InitFlags() {
	flag.StringVar(
		&c.brokers,
		"kafka-brokers",
		"localhost:9092",
		"Kafka broker addresses, comma-separated",
	)
}

func (c *producerComp) Activate(_ sctx.ServiceContext) error {
	c.logger = sctx.GlobalLogger().GetLogger(c.id)

	c.logger.Infof("Connecting to Kafka... %s", c.brokers)

	c.writer = &kafka.Writer{
		Addr:         kafka.TCP(strings.Split(c.brokers, ",")...),
		BatchTimeout: time.Millisecond * 100,
		RequiredAcks: kafka.RequireAll,
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	var conn *kafka.Conn
	var err error
	for _, broker := range strings.Split(c.brokers, ",") {
		conn, err = dialer.DialContext(context.Background(), "tcp", broker)
		if err != nil {
			break
		}

		conn.Close()
	}

	if err != nil {
		return err
	}

	c.logger.Info("Connected to Kafka")

	return nil
}

func (c *producerComp) Stop() error {
	c.logger.Info("Closing Kafka producer...")
	return c.writer.Close()
}

func (c *producerComp) PublishMessage(ctx context.Context, topic string, key []byte, value []byte) error {
	message := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	c.logger.Infof("Publishing message to Kafka with topic %s and message %+v\n", topic, message)

	err := c.writer.WriteMessages(ctx, message)
	if err != nil {
		c.logger.Error("Failed to produce message", err.Error())
		return err
	}

	return nil
}

func (c *producerComp) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to serialize message", err.Error())
		return err
	}

	return c.PublishMessage(ctx, topic, []byte(key), payload)
}
