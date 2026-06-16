package msk

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"sync"
	"time"
)

// MSKAsyncProducer represents a managed Kafka queue instance.
type asyncProducer struct {
	brokers    []string
	producer   sarama.AsyncProducer
	producerMu sync.Mutex
	configure  *sarama.Config
}

// NewMSKAsyncProducerQueue creates a new MSKAsyncProducerQueue instance.
func NewMSKAsyncProducerAsyncProducer(brokers []string, configure *sarama.Config) IProducer {
	return &asyncProducer{
		brokers:   brokers,
		configure: configure,
	}
}

// connectProducer connects to Kafka broker and creates a producer.
func (a *asyncProducer) InitMskProducer(ctx context.Context, client sarama.Client) error {
	newProducer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return err
	}

	a.producerMu.Lock()
	defer a.producerMu.Unlock()
	a.producer = newProducer
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for {
			select {
			case err = <-a.producer.Errors():
				fmt.Println(time.Now(), err)
			case <-signals:
				a.producer.AsyncClose() // Trigger a shutdown of the producer.
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

// SendMessage sends a message to the specified topic.
func (a *asyncProducer) SendMessage(msg *sarama.ProducerMessage) error {
	var err error
	if a.producer == nil {
		return fmt.Errorf("asyncProducer is nil")
	}

	select {
	case a.producer.Input() <- msg:
	case errs := <-a.producer.Errors():
		err = errs.Err
		break
	case <-time.After(time.Second * 1):
		break
	case <-a.producer.Successes():
		break
	}
	return err
}

func (a *asyncProducer) SendMessages(msgList []*sarama.ProducerMessage) error {
	return nil
}

func (a *asyncProducer) Close() {
	a.producerMu.Lock()
	defer a.producerMu.Unlock()
	a.producer.Close()
	return
}
