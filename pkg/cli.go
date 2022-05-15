package pkg

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Options struct {
	configFilename string
}

type Cli struct {
	Options *Options
}

func NewCli() *Cli {
	options := parseOptionsFromFlags()
	return &Cli{Options: options}
}

func parseOptionsFromFlags() *Options {
	configFilename := flag.String("f", "", "the path to configMap file")
	flag.Parse()
	if *configFilename == "" {
		flag.Usage()
		os.Exit(1)
	}
	return &Options{configFilename: *configFilename}
}

func (c *Cli) ReadEventorConfigs() *[]EventorConfig {
	yamlConfig, err := os.ReadFile(c.Options.configFilename)
	if err != nil {
		panic(err)
	}
	eventorConfigs := NewEventorConfigs(yaml.Unmarshal, yamlConfig)
	return &eventorConfigs
}
