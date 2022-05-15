package pkg

import (
	"context"
)

type EventHandler struct {
	config    *EventHandlerConfig
	logger    Logger
	apiCaller ApiCaller
}

type ApiCaller interface {
	Configure(logger Logger, config *EventHandlerConfig) error
	Call(ctx context.Context, httpMethod string, url string, body []byte, headers HttpHeaders) (responseBody []byte, httpStatusCode int, err error)
	Close() error
}

func NewEventHandler(logger Logger, config *EventHandlerConfig, apiCaller ApiCaller) (*EventHandler, error) {
	err := apiCaller.Configure(logger, config)
	if err != nil {
		return nil, err
	}
	return &EventHandler{config, logger, apiCaller}, nil
}

func (h *EventHandler) Handle(ctx context.Context, msg *Message) ([]byte, int, error) {
	h.logger.Debugf("calling endpoint: '%v %v'", h.config.Method, h.config.Url)
	return h.apiCaller.Call(ctx, h.config.Method, h.config.Url, msg.Value, h.config.Headers)
}

func (h *EventHandler) Close() {
	err := h.apiCaller.Close()
	if err != nil {
		h.logger.Errorf("error while closing apiCaller: %s", err)
		return
	}
}
