package pkg

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type KafkaMessageConsumer struct {
	Consumer    *kafka.Consumer
	logger      Logger
	config      *KafkaConsumerConfig
	lastMessage *kafka.Message
	topic       string
}

func (k *KafkaMessageConsumer) Configure(logger Logger, config *KafkaConsumerConfig, topic string) error {
	k.logger = logger
	k.config = config
	k.topic = topic
	if k.Consumer == nil {
		configMap, err := consumerConfigToConfigMap(config)
		if err != nil {
			return err
		}
		k.Consumer, err = kafka.NewConsumer(&configMap)
		if err != nil {
			logger.Errorf("error while creating consumer: %v", err)
			return err
		}
	}
	return nil
}

func (k *KafkaMessageConsumer) Subscribe() error {
	return k.Consumer.Subscribe(k.topic, nil)
}

func (k *KafkaMessageConsumer) Consume(ctx context.Context) (*Message, error) {
	run := true
	go func() {
		<-ctx.Done()
		run = false
	}()
	for run {
		msg, err := k.Consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			// Ignore any timeout errors (which will happen every 100ms)
			if err.(kafka.Error).Code() != kafka.ErrTimedOut {
				k.logger.Errorf("consumer error: %v (%v)", err, msg)
				return nil, err
			}
			continue
		}
		k.lastMessage = msg
		return &Message{
			Key:   msg.Key,
			Value: msg.Value,
		}, nil
	}
	return nil, fmt.Errorf("ctx is done but never received a message")
}

func (k *KafkaMessageConsumer) CommitLastMessage(ctx context.Context) error {
	_, err := k.Consumer.CommitMessage(k.lastMessage)
	if err == nil {
		k.lastMessage = nil
	}
	return err
}

func (k *KafkaMessageConsumer) Close() error {
	return k.Consumer.Close()
}

func consumerConfigToConfigMap(consumerConfig *KafkaConsumerConfig) (kafka.ConfigMap, error) {
	conf := kafka.ConfigMap{}
	for k, v := range *consumerConfig {
		if err := conf.Set(fmt.Sprintf("%v=%v", k, v)); err != nil {
			return nil, err
		}
	}
	return conf, nil
}
