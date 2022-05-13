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

func TestNewKafkaConfigFromYAML(t *testing.T) {
	yamlConfig := []byte(`
- name: MY_FIRST_TOPIC-to-MY_RESPONSE_TOPIC
  eventListener:
    type: kafka
    topic: MY_FIRST_TOPIC
    consumerConfig:
      bootstrap.servers: localhost, 127.0.0.1
      group.id: myGroup
      auto.offset.reset: earliest
  eventHandler:
    type: http
    url: http://localhost:8080/payload
    headers:
      Content-Type: application/json
      Authorization: Basic ...
  eventResultProcessor:
    type: kafka
    topic: MY_RESPONSE_TOPIC
    when:
      - onSuccess
    producerConfig:
      bootstrap.servers: localhost
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
	expectEqual(t, eventListener.ConsumerConfig["group.id"], "myGroup")
	expectEqual(t, eventListener.ConsumerConfig["auto.offset.reset"], "earliest")

	// check eventHandler
	eventHandler := firstEventorConfig.EventHandler
	expectEqual(t, eventHandler.Type, "http")
	expectEqual(t, eventHandler.Url, "http://localhost:8080/payload")
	expectEqual(t, eventHandler.Headers["Content-Type"], "application/json")
	expectEqual(t, eventHandler.Headers["Authorization"], "Basic ...")

	// check eventResultProcessor
	eventResultProcessor := firstEventorConfig.EventResultProcessor
	expectEqual(t, eventResultProcessor.Type, "kafka")
	expectEqual(t, eventResultProcessor.Topic, "MY_RESPONSE_TOPIC")
	expectEqual(t, eventResultProcessor.ProducerConfig["bootstrap.servers"], "localhost")
}
