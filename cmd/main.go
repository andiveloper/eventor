package main

import (
	"context"
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
	eventorConfigs := pkg.NewEventorConfigs(yaml.Unmarshal, yamlConfig)
	fmt.Println(eventorConfigs)

	ctx := context.Background()
	for _, config := range eventorConfigs {
		go pkg.NewEventor(config).Run(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("err: %s\n", err)
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
	}
}
