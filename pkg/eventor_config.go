package pkg

type EventorConfig struct {
	Name                 string
	EventListener        EventListener        `yaml:"eventListener"`
	EventHandler         EventHandler         `yaml:"eventHandler"`
	EventResultProcessor EventResultProcessor `yaml:"eventResultProcessor"`
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

// EventListener defines a Kafka consumer
type EventListener struct {
	Type           string // 'kafka' is the only supported type for now
	Topic          string
	ConsumerConfig map[string]string `yaml:"consumerConfig"`
}

// EventHandler defines an HTTP endpoint to which the events' payload is sent to
type EventHandler struct {
	Type    string // 'http' is the only supported type for now
	Url     string
	Headers map[string]string
}

// EventResultProcessor defines a Kafka producer which posts the response from the EventHandler to a Kafka topic
type EventResultProcessor struct {
	Type           string // 'kafka' is the only supported type for now
	Topic          string
	When           []string
	ProducerConfig map[string]string `yaml:"producerConfig"`
}
