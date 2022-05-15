package pkg

type EventorConfig struct {
	Name                 string
	LogLevel             int
	EventListener        EventListenerConfig        `yaml:"eventListener"`
	EventHandler         EventHandlerConfig         `yaml:"eventHandler"`
	EventResultProcessor EventResultProcessorConfig `yaml:"eventResultProcessor"`
}

type Unmarshaller interface {
	Unmarshal(in []byte, out interface{}) error
}

func NewEventorConfigs(unmarshaller func(in []byte, out interface{}) error, data []byte) []EventorConfig {
	var eventorConfigs []EventorConfig
	err := unmarshaller(data, &eventorConfigs)
	if err != nil {
		panic(err)
	}
	return eventorConfigs
}

type KafkaConsumerConfig = map[string]string

// EventListenerConfig defines a Kafka producer
type EventListenerConfig struct {
	Type           string // 'kafka' is the only supported type for now
	Topic          string
	ConsumerConfig KafkaConsumerConfig `yaml:"consumerConfig"`
}

type HttpHeaders = map[string]string

// EventHandlerConfig defines an HTTP endpoint to which the events' payload is sent to
type EventHandlerConfig struct {
	Type    string // 'http' is the only supported type for now
	Method  string
	Url     string
	Headers HttpHeaders
}

type KafkaProducerConfig = map[string]string

// EventResultProcessorConfig defines a Kafka producer which posts the response from the EventHandlerConfig to a Kafka topic
type EventResultProcessorConfig struct {
	Type           string // 'kafka' is the only supported type for now
	Topic          string
	When           []string
	ProducerConfig KafkaProducerConfig `yaml:"producerConfig"`
}
