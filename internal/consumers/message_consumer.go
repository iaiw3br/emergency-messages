// Package consumers
// Handles reading messages from Kafka queue
package consumers

import (
	"fmt"
	"log/slog"
)

type Consumer struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Consumer {
	return &Consumer{
		log: log,
	}
}

func (c *Consumer) Read(messageBytes []byte) {
	fmt.Println("READ CONSUMERS", string(messageBytes))
}
