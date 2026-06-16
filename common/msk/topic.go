package msk

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

func InitTopics(brokers []string, topics []string, configure *sarama.Config) error {
	log.Println("init topics...")
	client, err := sarama.NewClient(brokers, configure)
	if err != nil {
		return errors.New(fmt.Sprintf("InitTopics NewClient error: %s", err))
	}

	clusterAdmin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return errors.New(fmt.Sprintf("InitTopics NewClusterAdminFromClient error: %s", err))
	}

	for _, topic := range topics {
		err = clusterAdmin.CreateTopic(topic, &sarama.TopicDetail{NumPartitions: -1, ReplicationFactor: -1}, false)
		if err != nil && !errors.Is(err, sarama.ErrTopicAlreadyExists) {
			return errors.New(fmt.Sprintf("InitTopics CreateTopic error: %s", err))
		}
	}

	return nil
}
