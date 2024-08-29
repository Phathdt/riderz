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
	id            string
	brokers       string
	logger        sctx.Logger
	readers       map[string]*kafka.Reader
	batchReaders  map[string]*kafka.Reader
	handlers      map[string]kcomp.HandlerFunc
	batchHandlers map[string]kcomp.HandlerMultiFunc
	mu            sync.Mutex
	stopChan      chan struct{}
}

func New(id string) *consumerComp {
	return &consumerComp{
		id:            id,
		readers:       make(map[string]*kafka.Reader),
		batchReaders:  make(map[string]*kafka.Reader),
		handlers:      make(map[string]kcomp.HandlerFunc),
		batchHandlers: make(map[string]kcomp.HandlerMultiFunc),
		stopChan:      make(chan struct{}),
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
	c.logger.Infof("Connecting to Kafka... %s", c.brokers)

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

	for _, reader := range c.batchReaders {
		if err := reader.Close(); err != nil {
			c.logger.Error("Error closing reader", "error", err.Error())
		}
	}
	return nil
}

func (c *consumerComp) Subscribe(groupId, topic string, handlerFunc kcomp.HandlerFunc) {
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

func (c *consumerComp) BatchSubscribe(groupId string, topic string, batchTimeout time.Duration, batchSize int, handlerFunc kcomp.HandlerMultiFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.batchReaders[topic]; exists {
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

	c.batchReaders[topic] = reader
	c.batchHandlers[topic] = handlerFunc

	go c.consumeBatchMessages(topic, batchTimeout, batchSize)

	c.logger.Infof("Batch subscribed to topic %s with groupId %s", topic, groupId)
}

func (c *consumerComp) consumeBatchMessages(topic string, batchTimeout time.Duration, batchSize int) {
	reader := c.batchReaders[topic]
	handler := c.batchHandlers[topic]

	for {
		select {
		case <-c.stopChan:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), batchTimeout)
			messages := make([]*kcomp.Message, 0, batchSize)
			kafkaMessages := make([]kafka.Message, 0, batchSize)

			for i := 0; i < batchSize; i++ {
				m, err := reader.FetchMessage(ctx)
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						break // Timeout reached, process the batch
					}
					c.logger.Errorf("Error fetching message topic %s with error %s", topic, err.Error())
					continue
				}

				messages = append(messages, &kcomp.Message{
					Key:     m.Key,
					Payload: m.Value,
				})
				kafkaMessages = append(kafkaMessages, m)

				if len(messages) == batchSize {
					break // Batch size reached, process the batch
				}
			}

			cancel()

			if len(messages) > 0 {
				func() {
					defer func() {
						if r := recover(); r != nil {
							c.logger.Errorf("Panic occurred while handling batch topic %s with recover %+v\n", topic, r)
						}
					}()

					if err := handler(messages); err != nil {
						c.logger.Errorf("Error handling batch topic %s with error %s\n", topic, err.Error())
					}
				}()

				err := reader.CommitMessages(context.Background(), kafkaMessages...)
				if err != nil {
					c.logger.Errorf("Failed to commit messages topic %s with error %s\n", topic, err.Error())
				}
			}
		}
	}
}
