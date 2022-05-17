# eventor

`eventor` is an ultra lightweight proxy which turns your synchronous REST API into an asynchronous, event-driven
microservice. It does that by consuming an event from a Kafka topic, sending the event to your REST API endpoint and
publishing the response onto another Kafka topic.

In the cloud-native Kubernetes world a common use case could be to add `eventor` as a sidecar container to your existing pod.
That way your so far synchronous RESTful service is immediately opened up for an event-driven architecture.

![eventor](eventor.png?raw=true "eventor")

## Sample config:

```yaml
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
```

This config starts an event listener which consumes events from a Kafka topic called `MY_FIRST_TOPIC`. The consumed
event payload is then send to `http://localhost:8080/payload` for processing.

If the returned HTTP status code is <400 (because `when` is set to `onSuccess`) the response body is sent to
another Kafka topic called `MY_RESPONSE_TOPIC`.

Only if the call to the http endpoint AND sending the event to 'MY_RESPONSE_TOPIC' was successful the message is
committed.

## Additional librdkafka setup to build/test/run on MacOS/arm64

```bash
brew install openssl
brew install librdkafka
brew install pkg-config
export PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@3/lib/pkgconfig"
go build -tags dynamic -o eventor cmd/main.go
```

## Sample usage

```bash
# Start a single node Kafka cluster
docker-compose -f ./kafka/docker-compose.yaml up -d

# Open a new terminal and start an interactive console producer which can send events to "MY_FIRST_TOPIC":
docker exec -it kafka_kafka_1 /bin/sh -c "kafka-console-producer.sh --bootstrap-server localhost:9092 --topic MY_FIRST_TOPIC"

# Open a new terminal and start an interactive console consumer which reads events from the topic "MY_RESPONSE_TOPIC"
docker exec -it kafka_kafka_1 /bin/sh -c "kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic MY_RESPONSE_TOPIC"

# Open a new terminal and start the sample HTTP server which will print the received headers and body and send the body as-is back to the client
go run tools/http_server.go

# Open a new terminal and  start the eventor service with the sample config:
# 'eventor' will read any message from 'MY_FIRST_TOPIC', send the message to the local HTTP server, and publish the response to 'MY_RESPONSE_TOPIC'
# Note: '-tags dynamic' must be used on Mac/arm64 so that the correct librdkafka is used
go run -tags dynamic cmd/main.go -f sample_config.yaml
```

When all services are up and running you can now continue and send your first event:

```bash
# Sent an event using the interactive console producer by typing 'Hello World' and hitting enter:
$ docker exec -it kafka_kafka_1 /bin/sh -c "kafka-console-producer.sh --bootstrap-server localhost:9092 --topic MY_FIRST_TOPIC"
>Hello World!

# The eventor service will show:
$ go run -tags dynamic cmd/main.go -f sample_config.yaml
2022/05/16 15:42:12 DEBUG - handling message: MY_FIRST_TOPIC[0]@184
2022/05/16 15:42:12 DEBUG - calling endpoint: 'POST http://localhost:8080/payload'
2022/05/16 15:42:12 DEBUG - successfully produced EventHandler result: MY_RESPONSE_TOPIC[0]@97
2022/05/16 15:42:38 DEBUG - handling message: MY_FIRST_TOPIC[0]@185
2022/05/16 15:42:38 DEBUG - calling endpoint: 'POST http://localhost:8080/payload'
2022/05/16 15:42:38 DEBUG - successfully produced EventHandler result: MY_RESPONSE_TOPIC[0]@98

# The sample HTTP server will show:
$ go run tools/http_server.go
Listening on port 8080...
---------------------- request ----------------------
Accept-Encoding: gzip
User-Agent: Go-http-client/1.1
Content-Length: 12
Authorization: Basic ...
Content-Type: application/text

Hello World!


# The consumer will show:
$ docker exec -it kafka_kafka_1 /bin/sh -c "kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic MY_RESPONSE_TOPIC"
Hello World!
```

## Future work

- support for custom Kafka event keys (null is used for now, which does not guarantee event ordering when multiple
  partitions are used)
- support for proper error handling if the call to the `eventHandler` fails or responds with a none 2xx status code
- support for additional HTTP authentication schemes like OAuth2.0, ...
- support for other `eventListener`, `eventHandler` and `eventResultProcessor` types like NATS or RabbitMQ

## Useful commands

```bash
# build:
go build -tags dynamic -o eventor cmd/main.go

# test:
go test -tags dynamic ./pkg ./cmd

# run:
go run -tags dynamic cmd/main.go

# run sample:
go run -tags dynamic cmd/main.go -f sample_config.yaml

# start sample http server:
go run tools/http_server.go

# start single node kafka:
docker-compose -f ./kafka/docker-compose.yaml up -d

# start consumer:
docker exec -it kafka_kafka_1 /bin/sh -c "kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic MY_RESPONSE_TOPIC"

# start producer:
docker exec -it kafka_kafka_1 /bin/sh -c "kafka-console-producer.sh --bootstrap-server localhost:9092 --topic MY_FIRST_TOPIC"
```
