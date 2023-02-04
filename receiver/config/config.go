package config

import (
	"fmt"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/tarekbadrshalaan/goStuff/configuration"
)

var (
	readConfigOnce sync.Once
	internalConfig Config
)

// Config : application configuration
type Config struct {
	APIPort int `json:"RECEIVER_PORT_INTERNAL" envconfig:"RECEIVER_PORT_INTERNAL"`

	KafkaBrokers   string `json:"KAFKA_BROKERS" envconfig:"KAFKA_BROKERS"`
	KafkaOldest    bool   `json:"KAFKA_OLDEST_OFFSET" envconfig:"KAFKA_OLDEST_OFFSET"`
	KafkaGroupName string `json:"RECEIVER_KAFKA_GROUP" envconfig:"RECEIVER_KAFKA_GROUP"`

	RawURLTopic string `json:"RAW_URL_KAFKA_TOPIC" envconfig:"RAW_URL_KAFKA_TOPIC"`
}

// Configuration : get application configuration
func Configuration() Config {
	readConfigOnce.Do(func() {
		jsonerr := configuration.JSON("config.json", &internalConfig)
		if jsonerr != nil {
			fmt.Println(jsonerr)
			// get configuration from environment variables
			err := envconfig.Process("", &internalConfig)
			if err != nil {
				err = fmt.Errorf("Error while initiating app configuration : %v", err)
				panic(err)
			}
		}
	})
	return internalConfig
}
