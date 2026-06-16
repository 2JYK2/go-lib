package msk

import (
	"context"
	"github.com/IBM/sarama"
)

type IEventHandler interface {
	Process(msg sarama.ConsumerMessage)
	SetProducer(producer IProducer)
}

type handler struct {
	producer IProducer
	handler  IEventHandler
}

func (h handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.WithValue(session.Context(), "name", "consume claim")
	for {
		select {
		case <-ctx.Done():
			err := h.Cleanup(session)
			session.Commit()
			return err

		case message := <-claim.Messages():
			if h.handler != nil && message != nil {
				h.handler.Process(*message)
				session.MarkMessage(message, "")
			}
		}
	}
}
