package msk

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"strings"
	"sync"
)

type threadConsumer struct {
	ctx          context.Context
	wg           *sync.WaitGroup
	configure    *sarama.Config
	threadName   string
	brokers      []string
	groupId      string
	topic        string
	handler      IEventHandler
	needProducer bool
}

// InitMskThreadConsumer 初始化线程消费端
func InitMskThreadConsumer(ctx context.Context, wg *sync.WaitGroup, configure *sarama.Config, brokers []string, topic string, groupId string, handler IEventHandler, needProducer bool) *threadConsumer {
	threadName := fmt.Sprintf("topic:%s group:%s", topic, groupId)
	readerCtx := context.WithValue(ctx, "name", threadName)
	c := &threadConsumer{
		ctx:          readerCtx,
		wg:           wg,
		configure:    configure,
		threadName:   threadName,
		brokers:      brokers,
		groupId:      groupId,
		topic:        topic,
		handler:      handler,
		needProducer: needProducer,
	}
	return c
}

// Run 线程主逻辑，建议配合外层协程管理使用
func (t threadConsumer) Run() error {
	t.wg.Add(1)
	defer t.wg.Done()

	defer log.Println("consume quit", t.threadName)

	consumeLoopCtx := context.WithValue(t.ctx, "name", "Msk Consume Loop")
	if t.needProducer {
		//TODO
		/*mskObject := NewMSKSyncProducerAsyncProducer(t.brokers, t.configure)
		err := mskObject.InitMskProducer(t.ctx, client)
		if err != nil {
			return err
		}
		t.handler.SetProducer(mskObject)*/
	}

	consumer, err := sarama.NewConsumerGroup(t.brokers, t.groupId, t.configure)
	if err != nil {
		return errors.New(fmt.Sprintf("%s NewConsumerGroupFromClient error: %s", t.threadName, err))
	}

	defer func() {
		_ = consumer.Close()
	}()

	log.Println("consume from", t.threadName)
	for {
		select {
		case <-consumeLoopCtx.Done():
			log.Println("consume context canceled", t.threadName)
			return nil
		default:
			handlerCtx := context.WithValue(t.ctx, "name", "Msk Handler")
			topics := strings.Split(t.topic, ",")
			if err = consumer.Consume(handlerCtx, topics, handler{handler: t.handler}); err != nil {
				log.Println(errors.New(fmt.Sprintf("consume:%s error: %s", t.threadName, err)))
			}
			log.Println(errors.New(fmt.Sprintf("connecting  consume:%s error: %s", t.threadName, err)))
		}
	}
}
