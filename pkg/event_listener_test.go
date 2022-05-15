package pkg

import (
	"context"
	"testing"
)

func NewTestEventListener(t *testing.T) (*TestMessageConsumer, *EventListener) {
	testMessageConsumer := &TestMessageConsumer{}
	listener, err := NewEventListener(DefaultLogger(DEBUG), &EventListenerConfig{}, testMessageConsumer)
	if err != nil {
		t.Error(err)
	}
	return testMessageConsumer, listener
}

func TestNewEventListener(t *testing.T) {
	// when
	testMessageConsumer, _ := NewTestEventListener(t)

	// then
	if !testMessageConsumer.configureCalled {
		t.Errorf("messageConsumer.Configure was not called")
	}

	if !testMessageConsumer.subscribeCalled {
		t.Errorf("messageConsumer.Subscribe was not called")
	}
}

func TestEventListener_Listen(t *testing.T) {
	// given
	testMessageConsumer, listener := NewTestEventListener(t)
	testMessageConsumer.ConsumeMessage = &Message{
		Key:   nil,
		Value: nil,
	}
	handlerFuncCalled := false

	// when
	ctx, cancel := context.WithCancel(context.TODO())
	listener.Listen(ctx, func(ctx context.Context, msg *Message) error {
		handlerFuncCalled = true
		defer cancel()
		return nil
	})

	// then
	if !handlerFuncCalled {
		t.Errorf("handlerFuncCalled was not called")
	}

	if !testMessageConsumer.consumeCalled {
		t.Errorf("messageConsumer.Consume was not called")
	}

	if !testMessageConsumer.commitLastMessageCalled {
		t.Errorf("messageConsumer.CommitLastMessage was not called")
	}
}

func TestEventListener_Close(t *testing.T) {
	// given
	testMessageConsumer, listener := NewTestEventListener(t)

	// when
	listener.Close()

	// then
	if !testMessageConsumer.closeCalled {
		t.Errorf("messageConsumer.Close was not called")
	}
}

type TestMessageConsumer struct {
	configureCalled         bool
	subscribeCalled         bool
	consumeCalled           bool
	ConsumeMessage          *Message
	commitLastMessageCalled bool
	closeCalled             bool
}

func (t *TestMessageConsumer) Configure(logger Logger, config *KafkaConsumerConfig, topic string) error {
	t.configureCalled = true
	return nil
}

func (t *TestMessageConsumer) Subscribe() error {
	t.subscribeCalled = true
	return nil
}

func (t *TestMessageConsumer) Consume(ctx context.Context) (*Message, error) {
	t.consumeCalled = true
	return t.ConsumeMessage, nil
}

func (t *TestMessageConsumer) CommitLastMessage(ctx context.Context) error {
	t.commitLastMessageCalled = true
	return nil
}

func (t *TestMessageConsumer) Close() error {
	t.closeCalled = true
	return nil
}
