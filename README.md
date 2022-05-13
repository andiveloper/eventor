# eventor

Ultra lightweight service to turn your synchronous REST API into an asynchronous, event-driven microservice

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