package pkg

import (
	"context"
)

type EventListener struct {
	config   *EventListenerConfig
	consumer MessageConsumer
	logger   Logger
}

type MessageConsumer interface {
	Configure(logger Logger, config *KafkaConsumerConfig, topic string) error
	Subscribe() error
	Consume(ctx context.Context) (*Message, error)
	CommitLastMessage(ctx context.Context) error
	Close() error
}

func NewEventListener(logger Logger, config *EventListenerConfig, consumer MessageConsumer) (*EventListener, error) {
	err := consumer.Configure(logger, &config.ConsumerConfig, config.Topic)
	if err != nil {
		return nil, err
	}
	if err := consumer.Subscribe(); err != nil {
		logger.Errorf("error: %v", err)
		return nil, err
	}
	return &EventListener{config, consumer, logger}, nil
}

func (l *EventListener) Listen(ctx context.Context, handleMessage func(ctx context.Context, msg *Message) error) {
	run := true
	go func() {
		<-ctx.Done()
		run = false
	}()

	for run {
		msg, err := l.consumer.Consume(ctx)
		if err != nil {
			// Ignore any errors and assume the consumer will recover
			continue
		}
		// handle message
		l.logger.Debugf("handling message: %v", msg)
		if err := handleMessage(ctx, msg); err != nil {
			l.logger.Errorf("error while handling message on topic %v: %v (%v)", l.config.Topic, err, msg)
		}
		// commit the message
		err = l.consumer.CommitLastMessage(ctx)
		if err != nil {
			l.logger.Errorf("error while committing message on topic %v: %v (%v)", l.config.Topic, err, msg)
		}
	}
}

func (l *EventListener) Close() {
	err := l.consumer.Close()
	if err != nil {
		l.logger.Errorf("error while closing consumer: %s", err)
		return
	}
}
