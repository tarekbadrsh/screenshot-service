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
	Rresolution string `json:"RESOLUTION" envconfig:"RESOLUTION"`

	ChromeTimeOut          int    `json:"CHROME_TIMEOUT" envconfig:"CHROME_TIMEOUT"`
	ChromeTimeBudget       int    `json:"CHROME_TIME_BUDGET" envconfig:"CHROME_TIME_BUDGET"`
	ScrapySplashHost       string `json:"SCRAPY_SPLASH_HOST" envconfig:"SCRAPY_SPLASH_HOST"`
	ParallerExecutionCount int    `json:"PARALLER_EXECUTION_COUNT" envconfig:"PARALLER_EXECUTION_COUNT"`
	RetryCount             int    `json:"RETRY_COUNT" envconfig:"RETRY_COUNT"`

	StoragePath         string `json:"STORAGE_PATH" envconfig:"STORAGE_PATH"`
	KafkaClusterVersion string `json:"KAFKA_CLUSTER_VERSION" envconfig:"KAFKA_CLUSTER_VERSION"`
	KafkaBrokers        string `json:"KAFKA_BROKERS" envconfig:"KAFKA_BROKERS"`
	KafkaOldest         bool   `json:"KAFKA_OLDEST_OFFSET" envconfig:"KAFKA_OLDEST_OFFSET"`
	KafkaGroupName      string `json:"SCREEN_SHOT_KAFKA_GROUP" envconfig:"SCREEN_SHOT_KAFKA_GROUP"`

	RawURLTopic      string `json:"RAW_URL_KAFKA_TOPIC" envconfig:"RAW_URL_KAFKA_TOPIC"`
	ResultServiceURL string `json:"RESULT_SERVICE_URL" envconfig:"RESULT_SERVICE_URL"`
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
