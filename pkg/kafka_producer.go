package pkg

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaMessageProducer struct {
	Producer *kafka.Producer
	logger   Logger
	config   *KafkaProducerConfig
	topic    string
}

func (k *KafkaMessageProducer) Configure(logger Logger, config *KafkaProducerConfig, topic string) error {
	k.logger = logger
	k.config = config
	k.topic = topic
	if k.Producer == nil {
		configMap, err := producerConfigToConfigMap(k.config)
		if err != nil {
			return err
		}
		k.Producer, err = kafka.NewProducer(&configMap)
		if err != nil {
			k.logger.Errorf("error: %v", err)
			return err
		}
	}
	return nil
}

func (k *KafkaMessageProducer) Produce(ctx context.Context, msg *Message) error {
	run := true
	go func() {
		<-ctx.Done()
		run = false
	}()
	deliveryChan := make(chan kafka.Event, 1)
	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Key:            msg.Key,
		Value:          msg.Value,
	}, deliveryChan)
	if err != nil {
		k.logger.Errorf("producer error on topic %v: %v", k.topic, err)
	}
	// wait for message delivery
	messagesThatNeedToBeDelivered := 1
	for messagesThatNeedToBeDelivered > 0 && run {
		messagesThatNeedToBeDelivered = k.Producer.Flush(100)
		evt := <-deliveryChan
		msg := evt.(*kafka.Message)
		k.logger.Debugf("successfully produced EventHandler result: %v", msg)
	}
	return nil
}

func (k *KafkaMessageProducer) Close() error {
	k.Producer.Close()
	return nil
}

func producerConfigToConfigMap(producerConfig *KafkaProducerConfig) (kafka.ConfigMap, error) {
	conf := kafka.ConfigMap{}
	for k, v := range *producerConfig {
		if err := conf.Set(fmt.Sprintf("%v=%v", k, v)); err != nil {
			return nil, err
		}
	}
	return conf, nil
}
