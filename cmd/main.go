package main

import (
	"flag"
	"fmt"
	"github.com/andiveloper/eventor/pkg"
	"gopkg.in/yaml.v3"
	"os"
)

type Options struct {
	configFilename string
}

func ParseOptionsFromFlags() *Options {
	configFilename := flag.String("f", "", "the path to config file")
	flag.Parse()
	if *configFilename == "" {
		flag.Usage()
		os.Exit(1)
	}
	return &Options{configFilename: *configFilename}
}

func main() {
	options := ParseOptionsFromFlags()
	yamlConfig, err := os.ReadFile(options.configFilename)
	if err != nil {
		panic(err)
	}
	kafkaConfig := pkg.NewEventorConfigs(yaml.Unmarshal, yamlConfig)
	fmt.Println(kafkaConfig)
}
