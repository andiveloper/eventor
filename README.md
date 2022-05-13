# eventor

Eventor is an ultra lightweight proxy which turns your synchronous REST API into an asynchronous, event-driven
microservice. It does that by reading events from a Kafka topic, sending the events to your REST API endpoint and
posting the result onto another Kafka topic.

## Sample config:

```yaml
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
```

This config starts an event listener which consumes events from a Kafka topic called `MY_FIRST_TOPIC`. The consumed
event payload is then send to `http://localhost:8080/payload` for processing. If the returned HTTP status code is 2xx (
because `when` is set to `onSuccess`) the result/the HTTP body is sent to another Kafka topic called `MY_RESPONSE_TOPIC`
.

## Usage

```bash

```

## Future work

- support for other `eventListener`, `eventHandler` and `eventResultProcessor` types like NATS or RabbitMQ

