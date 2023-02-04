package messaging

import (
	"errors"
	"receiver/logger"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	// KafkaReadiness : check if kafka is ready
	KafkaReadiness = false

	initKafkaOnce sync.Once
	//go:generate mockgen -destination=../mocks/mock_SyncProducer.go -package=mocks github.com/Shopify/sarama SyncProducer
	kafkaProducer sarama.SyncProducer
)

func newKafkaProducerConfiguration() *sarama.Config {
	sarama.Logger = *logger.GetLogger()
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.ChannelBufferSize = 1
	return conf
}

// NewSaramaSyncProducer :
func NewSaramaSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	if len(brokers) == 0 {
		err := errors.New("No KAFKA BROKERS")
		logger.Error(err)
		return nil, err
	}
	kafka, err := sarama.NewSyncProducer(brokers, newKafkaProducerConfiguration())
	if err != nil {
		return nil, err
	}
	KafkaReadiness = true
	return kafka, nil
}

// Produce : Produce message to Kafka
func Produce(topic string, msg []byte) error {
	logger.Infof("message to produce %v", string(msg))
	partition, offset, err := kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("Produced at partion: %v, offset: %v", partition, offset)
	return nil
}

// InitializeKafka : initialize kafka clinet
func InitializeKafka(s sarama.SyncProducer) {
	initKafkaOnce.Do(func() {
		kafkaProducer = s
		logger.Info("Kafka has been initialized")
	})
}
