package data

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
)

// Producer is a type representing a Kafka producer.
type Producer struct {
	producer *kafka.Producer
	log      *slog.Logger
}

// NewProducer creates a new kafka producer
func NewProducer(brokerAddr string, log *slog.Logger) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerAddr,
	})
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: producer,
		log:      log,
	}, nil
}

// Send sends a message to the kafka topic
func (p *Producer) Send(messageBytes []byte) error {
	return p.send(EventTypeSend, messageBytes)
}

// Close closes the producer
func (p *Producer) Close() {
	p.producer.Close()
}

// send sends a message to the kafka topic
func (p *Producer) send(eventType EventType, value []byte) error {
	go func() {
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					p.log.With(slog.Any("message", string(value))).
						Error("Delivery failed", ev.TopicPartition.Error)
				} else {
					p.log.With(
						slog.Any("message", string(value)),
						slog.Any("topic", ev.TopicPartition)).
						Debug("Delivered message")
				}
			}
		}
	}()

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     (*string)(&eventType),
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
