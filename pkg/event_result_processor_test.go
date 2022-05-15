package pkg

import (
	"context"
	"testing"
)

func NewTestEventResultProcessor(t *testing.T) (*TestMessageProducer, *EventResultProcessor) {
	testMessageProducer := &TestMessageProducer{}
	listener, err := NewEventResultProcessor(DefaultLogger(DEBUG), &EventResultProcessorConfig{}, testMessageProducer)
	if err != nil {
		t.Error(err)
	}
	return testMessageProducer, listener
}

func TestNewEventResultProcessor(t *testing.T) {
	// given, when
	testMessageProducer, _ := NewTestEventResultProcessor(t)

	// then
	if !testMessageProducer.configureCalled {
		t.Errorf("messageProducer.Configure was not called")
	}
}

func TestEventResultProcessor_Produce(t *testing.T) {
	// given
	testMessageProducer, producer := NewTestEventResultProcessor(t)
	testMessage := &Message{
		Key:   nil,
		Value: nil,
	}

	// when
	ctx, cancel := context.WithCancel(context.TODO())
	err := producer.Produce(ctx, testMessage)
	if err != nil {
		t.Error(err)
	}
	cancel()

	// then
	if !testMessageProducer.produceCalled {
		t.Errorf("messageProducer.Produce was not called")
	}
}

func TestEventResultProcessor_Close(t *testing.T) {
	// given
	testMessageProducer, producer := NewTestEventResultProcessor(t)

	// when
	producer.Close()

	// then
	if !testMessageProducer.closeCalled {
		t.Errorf("messageProducer.Close was not called")
	}
}

type TestMessageProducer struct {
	configureCalled bool
	produceCalled   bool
	closeCalled     bool
}

func (t *TestMessageProducer) Configure(logger Logger, config *KafkaConsumerConfig, topic string) error {
	t.configureCalled = true
	return nil
}

func (t *TestMessageProducer) Produce(ctx context.Context, msg *Message) error {
	t.produceCalled = true
	return nil
}

func (t *TestMessageProducer) Close() error {
	t.closeCalled = true
	return nil
}
