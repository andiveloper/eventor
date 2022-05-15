package pkg

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func expectEqual[T comparable](t *testing.T, actual, expected T) {
	if actual != expected {
		t.Errorf(`Expected %v but got %v`, expected, actual)
	}
}

func TestNewEventorConfigs(t *testing.T) {
	yamlConfig := []byte(`
- name: MY_FIRST_TOPIC-to-MY_RESPONSE_TOPIC
  logLevel: 0 # 0 = DEBUG, 1 = INFO, 2 = WARN, 3 = ERROR
  eventListener:
    type: kafka
    topic: MY_FIRST_TOPIC # the topic from which events are consumed
    consumerConfig:
      bootstrap.servers: localhost, 127.0.0.1
      group.id: myVeryFirstGroup
      auto.offset.reset: earliest
      enable.auto.commit: false # auto commit should be disabled since the listener commits messages only if everything was fine
  eventHandler:
    type: http
    method: POST # the http method that is used for the request
    url: http://localhost:8080/payload # the URL to which the event is sent
    headers:
      Content-Type: application/text # defaults to application/text
      Authorization: Basic ...
  eventResultProcessor:
    type: kafka
    topic: MY_RESPONSE_TOPIC # the topic to which the response body received from the eventHandler is sent to
    when:
      - onSuccess # can be 'onSuccess' which means the response is published if HTTP statusCode < 400, or 'onError' for HTTP statusCode >= 400
    producerConfig:
      bootstrap.servers: localhost
      acks: all
`)

	eventorConfigs := NewEventorConfigs(yaml.Unmarshal, yamlConfig)
	j, _ := json.MarshalIndent(eventorConfigs, "", "    ")
	fmt.Printf("%v", string(j))

	expectEqual(t, len(eventorConfigs), 1)

	// check firstEventorConfig
	firstEventorConfig := eventorConfigs[0]
	expectEqual(t, firstEventorConfig.Name, "MY_FIRST_TOPIC-to-MY_RESPONSE_TOPIC")

	// check eventListener
	eventListener := firstEventorConfig.EventListener
	expectEqual(t, eventListener.Type, "kafka")
	expectEqual(t, eventListener.Topic, "MY_FIRST_TOPIC")
	expectEqual(t, eventListener.ConsumerConfig["bootstrap.servers"], "localhost, 127.0.0.1")
	expectEqual(t, eventListener.ConsumerConfig["group.id"], "myVeryFirstGroup")
	expectEqual(t, eventListener.ConsumerConfig["auto.offset.reset"], "earliest")

	// check eventHandler
	eventHandler := firstEventorConfig.EventHandler
	expectEqual(t, eventHandler.Type, "http")
	expectEqual(t, eventHandler.Url, "http://localhost:8080/payload")
	expectEqual(t, eventHandler.Headers["Content-Type"], "application/text")
	expectEqual(t, eventHandler.Headers["Authorization"], "Basic ...")

	// check eventResultProcessor
	eventResultProcessor := firstEventorConfig.EventResultProcessor
	expectEqual(t, eventResultProcessor.Type, "kafka")
	expectEqual(t, eventResultProcessor.Topic, "MY_RESPONSE_TOPIC")
	expectEqual(t, eventResultProcessor.ProducerConfig["bootstrap.servers"], "localhost")
}
