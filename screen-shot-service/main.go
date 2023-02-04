package main

import (
	"os"
	"os/signal"
	"screen-shot-service/config"
	"screen-shot-service/generator"
	"screen-shot-service/generator/scrapysplash"
	"screen-shot-service/logger"
	"screen-shot-service/messaging"
	"screen-shot-service/storage"
	"syscall"
)

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* logger initialize start */
	mylogger := logger.NewZapLogger()
	logger.InitializeLogger(&mylogger)
	defer logger.Close()
	/* logger initialize end */

	/* Generator Initialize start */
	// gen := chrome.NewChromeGenerator(c) // Google Chrome Client
	gen := scrapysplash.NewSplashGenerator(c) // Scrapy Splash Client
	generator.InitializeGenerator(&gen)
	/* Generator Initialize end */

	/* Result Service Initialize start */
	storage.InitializeResultService(c)
	/* Result Service Initialize end */

	/* Running Kafak start */
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	consumer := messaging.SaramaConsumer{
		Brokers:   []string{c.KafkaBrokers},
		Version:   c.KafkaClusterVersion,
		GroupName: c.KafkaGroupName,
		Topics:    []string{c.RawURLTopic},
		Oldest:    c.KafkaOldest,

		Ready:   make(chan bool, 0),
		Queue:   make(chan bool, c.ParallerExecutionCount),
		Handler: generator.GetGeneratorHandler(),
	}
	consumer.RunConsumerGroup(&sigterm)
	/* Running Kafak end */
}
