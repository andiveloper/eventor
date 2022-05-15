package main

import (
	"context"
	"github.com/andiveloper/eventor/pkg"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	// catch SIGINT
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	cli := pkg.NewCli()
	for _, config := range *cli.ReadEventorConfigs() {
		config := config
		go func() {
			logger := pkg.DefaultLogger(pkg.Level(config.LogLevel))
			err := pkg.NewEventor(config, &pkg.KafkaMessageConsumer{}, &pkg.HttpApiCaller{}, &pkg.KafkaMessageProducer{}).Run(ctx, logger)
			if err != nil {
				logger.Errorf("An error occurred: %v", err)
				cancel()
			}
		}()
	}

	logger := pkg.DefaultLogger(pkg.INFO)
	func() {
		select {
		case <-c:
			logger.Info("cancelled")
			cancel()
		case <-ctx.Done():
			logger.Info("done")
		}
	}()
}
