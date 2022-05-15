package pkg

import (
	"context"
	"golang.org/x/exp/slices"
	"testing"
)

func NewTestEventHandler(t *testing.T) (*TestApiCaller, *EventHandler) {
	testApiCaller := &TestApiCaller{}
	handler, err := NewEventHandler(DefaultLogger(DEBUG), &EventHandlerConfig{}, testApiCaller)
	if err != nil {
		t.Error(err)
	}
	return testApiCaller, handler
}

func TestNewEventHandler(t *testing.T) {
	// when
	testApiCaller, _ := NewTestEventHandler(t)

	// then
	if !testApiCaller.configureCalled {
		t.Errorf("ApiCaller.Configure was not called")
	}
}

func TestEventHandler_Handle(t *testing.T) {
	// given
	testApiCaller, handler := NewTestEventHandler(t)
	expectedResponse := []byte("response")
	testApiCaller.CallResponseBody = expectedResponse
	expectedStatusCode := 200
	testApiCaller.CallHttpStatusCode = expectedStatusCode

	// when
	response, statusCode, err := handler.Handle(context.TODO(), &Message{})
	if err != nil {
		return
	}

	// then
	if !testApiCaller.callCalled {
		t.Errorf("ApiCaller.callCalled was not called")
	}
	if slices.Compare(response, expectedResponse) != 0 {
		t.Errorf("response is not expectedResponse")
	}
	if statusCode != expectedStatusCode {
		t.Errorf("statusCode is not expectedStatusCode")
	}
}

func TestEventHandler_Close(t *testing.T) {
	// given
	testApiCaller, handler := NewTestEventHandler(t)

	// when
	handler.Close()

	// then
	if !testApiCaller.closeCalled {
		t.Errorf("ApiCaller.Close was not called")
	}
}

type TestApiCaller struct {
	configureCalled    bool
	callCalled         bool
	CallResponseBody   []byte
	CallHttpStatusCode int
	closeCalled        bool
}

func (t *TestApiCaller) Configure(logger Logger, config *EventHandlerConfig) error {
	t.configureCalled = true
	return nil
}

func (t *TestApiCaller) Call(ctx context.Context, httpMethod string, url string, body []byte, headers HttpHeaders) (responseBody []byte, httpStatusCode int, err error) {
	t.callCalled = true
	return t.CallResponseBody, t.CallHttpStatusCode, nil
}

func (t *TestApiCaller) Close() error {
	t.closeCalled = true
	return nil
}
