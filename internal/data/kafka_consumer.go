package data

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"log/slog"
	"projects/emergency-messages/internal/consumers"
	"time"
)

type Consumer struct {
	consumer        *kafka.Consumer
	messageConsumer *consumers.Consumer
	log             *slog.Logger
	done            chan bool
}

func NewConsumer(brokerAddr string, messageConsumer *consumers.Consumer, log *slog.Logger) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerAddr,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &Consumer{
		consumer:        consumer,
		messageConsumer: messageConsumer,
		log:             log,
		done:            make(chan bool),
	}, nil
}

func (c *Consumer) Close() {
	close(c.done)
	c.consumer.Close()
}

func (c *Consumer) Read() {
	events := getEvents()
	if err := c.consumer.SubscribeTopics(events, nil); err != nil {
		log.Fatal(err)
	}

	run := true
	for run {
		select {
		case <-c.done:
			run = false
		default:
			msg, err := c.consumer.ReadMessage(time.Second)
			if err == nil {
				switch *msg.TopicPartition.Topic {
				case EventTypeSend:
					c.messageConsumer.Read(msg.Value)
				}
				c.log.Info("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else {
				if !err.(kafka.Error).IsTimeout() {
					c.log.Info("Consumer error: %v (%v)\n", err, msg)
				}
			}
		}
	}
}
