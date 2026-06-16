package msk

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

type syncProducer struct {
	brokers    []string
	producer   sarama.SyncProducer
	producerMu sync.Mutex
	configure  *sarama.Config
}

type Message struct {
	Topic string
	Data  string
}

func NewMSKSyncProducerAsyncProducer(brokers []string, configure *sarama.Config) IProducer {
	return &syncProducer{
		brokers:   brokers,
		configure: configure,
	}
}

func (a *syncProducer) InitMskProducer(ctx context.Context, client sarama.Client) error {
	sp, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return errors.New(fmt.Sprintf("InitMskProducer NewSyncProducerFromClient error: %s", err))
	}

	a.producer = sp
	return nil
}

func (p *syncProducer) SendMessage(msg *sarama.ProducerMessage) error {
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		return errors.New(fmt.Sprintf("SendMessage error: %s", err))
	}
	return nil
}

func (p *syncProducer) SendMessages(msgList []*sarama.ProducerMessage) error {
	err := p.producer.SendMessages(msgList)
	if err != nil {
		return errors.New(fmt.Sprintf("SendMessages error: %s", err))
	}
	return nil
}

func (p *syncProducer) Close() {
	p.producerMu.Lock()
	defer p.producerMu.Unlock()
	p.producer.Close()
}
