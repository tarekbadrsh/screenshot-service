package messaging

import (
	"context"
	"os"
	"screen-shot-service/logger"

	"github.com/Shopify/sarama"
)

// SaramaConsumer : Saram consumer group handler implementation of interface : sarama.ConsumerGroupHandler
type SaramaConsumer struct {
	Brokers   []string
	Version   string //Kafka cluster version
	GroupName string
	Topics    []string
	Oldest    bool // consume from first offset

	Ready   chan bool
	Queue   chan bool
	Handler func([]byte) error
}

// RunConsumerGroup :
func (sc *SaramaConsumer) RunConsumerGroup(sigterm *chan os.Signal) {
	logger.Info("Starting a new kafka consumer")
	config := sc.initSarmaConfig()

	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(sc.Brokers, sc.GroupName, config)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		err := client.Consume(ctx, sc.Topics, sc)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	<-sc.Ready // Await till the consumer has been set up
	logger.Info("kafka consumer up and running!...")
	<-*sigterm // Await a sigterm signal before safely closing the consumer

	err = client.Close()
	if err != nil {
		logger.Error(err)
	}
}

func (sc *SaramaConsumer) initSarmaConfig() *sarama.Config {
	sarama.Logger = *logger.GetLogger()
	version, err := sarama.ParseKafkaVersion(sc.Version)
	if err != nil {
		logger.Fatal(err)
	}
	// New kafka config
	config := sarama.NewConfig()
	config.Version = version
	config.ClientID = sc.GroupName
	if sc.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	return config
}

// Setup : runs at the beginning of a new session, before ConsumeClaim
func (sc *SaramaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(sc.Ready)
	return nil
}

// Cleanup : runs at the end of a session, once all ConsumeClaim goroutines have exited
func (sc *SaramaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim : start a consumer loop of ConsumerGroupClaim's Messages().
func (sc *SaramaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		sc.Queue <- true
		go func(msg *sarama.ConsumerMessage) {
			logger.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", msg.Value, msg.Timestamp, msg.Topic)
			err := sc.Handler(msg.Value)
			if err != nil {
				logger.Error(err)
			}
			session.MarkMessage(msg, "")
			<-sc.Queue
		}(message)
	}
	return nil
}
