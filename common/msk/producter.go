package msk

import (
	"context"
	"github.com/IBM/sarama"
)

type IProducer interface {
	SendMessage(msg *sarama.ProducerMessage) error
	SendMessages(msgList []*sarama.ProducerMessage) error
	Close()
	InitMskProducer(ctx context.Context, client sarama.Client) error
}
