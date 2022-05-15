package pkg

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

type HttpApiCaller struct {
	logger Logger
	config *EventHandlerConfig
}

func (c *HttpApiCaller) Configure(logger Logger, config *EventHandlerConfig) error {
	c.logger = logger
	c.config = config
	return nil
}

func (c *HttpApiCaller) Call(ctx context.Context, httpMethod string, url string, body []byte, headers HttpHeaders) (responseBody []byte, httpStatusCode int, err error) {
	response, err := callApiInternal(ctx, httpMethod, url, body, headers)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			c.logger.Errorf("error while closing body: %v", err)
		}
	}()
	if err != nil {
		c.logger.Errorf("error during request: %v", err)
		return nil, 0, err
	}
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		c.logger.Errorf("error while reading response data: %v", err)
		return nil, 0, err
	}
	return responseBody, response.StatusCode, nil
}

func callApiInternal(ctx context.Context, httpMethod string, url string, body []byte, headers HttpHeaders) (*http.Response, error) {
	contentType, found := headers["Content-Type"]
	if !found {
		contentType = "application/text"
	}
	request, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	return client.Do(request)
}

func (c *HttpApiCaller) Close() error {
	return nil
}
