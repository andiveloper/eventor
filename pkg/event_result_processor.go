package pkg

import (
	"context"
)

type EventResultProcessor struct {
	config   *EventResultProcessorConfig
	producer MessageProducer
	logger   Logger
}

type MessageProducer interface {
	Configure(logger Logger, config *KafkaProducerConfig, topic string) error
	Produce(ctx context.Context, msg *Message) error
	Close() error
}

func NewEventResultProcessor(logger Logger, config *EventResultProcessorConfig, producer MessageProducer) (*EventResultProcessor, error) {
	err := producer.Configure(logger, &config.ProducerConfig, config.Topic)
	if err != nil {
		return nil, err
	}
	return &EventResultProcessor{config, producer, logger}, nil
}

func (r *EventResultProcessor) Produce(ctx context.Context, message *Message) error {
	if err := r.producer.Produce(ctx, message); err != nil {
		return err
	}
	return nil
}

func (r *EventResultProcessor) Close() {
	err := r.producer.Close()
	if err != nil {
		r.logger.Errorf("error while closing producer: %s", err)
		return
	}
}
