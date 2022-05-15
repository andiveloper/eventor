package pkg

import (
	"context"
	"golang.org/x/exp/slices"
)

type Eventor struct {
	config    *EventorConfig
	consumer  MessageConsumer
	apiCaller ApiCaller
	producer  MessageProducer
}

type Message struct {
	Key   []byte
	Value []byte
}

func NewEventor(config EventorConfig, consumer MessageConsumer, apiCaller ApiCaller, producer MessageProducer) *Eventor {
	return &Eventor{&config, consumer, apiCaller, producer}
}

func (e *Eventor) Run(ctx context.Context, logger Logger) error {
	eventListener, eventHandler, eventResultProcessor, err := e.setup(logger)
	if err != nil {
		return err
	}
	eventListener.Listen(ctx, func(ctx context.Context, msg *Message) error {
		return e.handleAndProcess(ctx, msg, eventHandler, eventResultProcessor)
	})
	return nil
}

func (e *Eventor) setup(logger Logger) (*EventListener, *EventHandler, *EventResultProcessor, error) {
	eventListener, err := NewEventListener(logger, &e.config.EventListener, e.consumer)
	if err != nil {
		return nil, nil, nil, err
	}

	eventHandler, err := NewEventHandler(logger, &e.config.EventHandler, e.apiCaller)
	if err != nil {
		return nil, nil, nil, err
	}

	eventResultProcessor, err := NewEventResultProcessor(logger, &e.config.EventResultProcessor, e.producer)
	if err != nil {
		return nil, nil, nil, err
	}
	return eventListener, eventHandler, eventResultProcessor, nil
}

func (e *Eventor) handleAndProcess(ctx context.Context, msg *Message, eventHandler *EventHandler, eventResultProcessor *EventResultProcessor) error {
	// TODO add retry logic with sleep, ...?
	response, statusCode, err := eventHandler.Handle(ctx, msg)
	if err != nil {
		return err
	}

	when := e.config.EventResultProcessor.When
	processResult := len(when) == 0 || (statusCode >= 400 && slices.Contains(when, "onError")) || (statusCode < 400 && slices.Contains(when, "onSuccess"))
	if processResult {
		err = eventResultProcessor.Produce(ctx, &Message{Key: nil, Value: response})
		if err != nil {
			return err
		}
	}
	return nil
}
