package kconsumercomp

import (
	"context"
	"errors"
	"flag"
	sctx "github.com/phathdt/service-context"
	"github.com/segmentio/kafka-go"
	"riderz/plugins/kcomp"
	"strings"
	"sync"
	"time"
)

type consumerComp struct {
	id       string
	brokers  string
	logger   sctx.Logger
	readers  map[string]*kafka.Reader
	handlers map[string]func(msg *kcomp.Message) error
	mu       sync.Mutex
	stopChan chan struct{}
}

func New(id string) *consumerComp {
	return &consumerComp{
		id:       id,
		readers:  make(map[string]*kafka.Reader),
		handlers: make(map[string]func(msg *kcomp.Message) error),
		stopChan: make(chan struct{}),
	}
}

func (c *consumerComp) ID() string {
	return c.id
}

func (c *consumerComp) InitFlags() {
	flag.StringVar(&c.brokers, "kafka-brokers", "localhost:9092", "Kafka broker addresses, comma-separated")
}

func (c *consumerComp) Activate(_ sctx.ServiceContext) error {
	c.logger = sctx.GlobalLogger().GetLogger(c.id)
	c.logger.Info("Kafka consumer component activated")
	return nil
}

func (c *consumerComp) Stop() error {
	c.logger.Info("Closing Kafka consumers...")
	close(c.stopChan)

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			c.logger.Error("Error closing reader", "error", err.Error())
		}
	}
	return nil
}

func (c *consumerComp) Subscribe(groupId, topic string, handlerFunc func(msg *kcomp.Message) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.readers[topic]; exists {
		c.logger.Warn("Already subscribed to topic", "topic", topic)
		return
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     strings.Split(c.brokers, ","),
		GroupID:     groupId,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
		MaxWait:     500 * time.Millisecond,
	})

	c.readers[topic] = reader
	c.handlers[topic] = handlerFunc

	go c.consumeMessages(topic)

	c.logger.Infof("Subscribed to topic %s with groupId %s", topic, groupId)
}

func (c *consumerComp) consumeMessages(topic string) {
	reader := c.readers[topic]
	handler := c.handlers[topic]

	for {
		select {
		case <-c.stopChan:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			m, err := reader.FetchMessage(ctx)
			cancel()

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					continue // This is normal, just retry
				}
				c.logger.Error("Error fetching message", "topic", topic, "error", err.Error())
				continue
			}

			msg := &kcomp.Message{
				Key:     m.Key,
				Payload: m.Value,
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						c.logger.Error("Panic occurred while handling message", "topic", topic, "recover", r)
					}
				}()

				if err = handler(msg); err != nil {
					c.logger.Error("Error handling message", "topic", topic, "error", err.Error())
				}
			}()

			err = reader.CommitMessages(context.Background(), m)
			if err != nil {
				c.logger.Error("Failed to commit message", "topic", topic, "error", err.Error())
			}
		}
	}
}
