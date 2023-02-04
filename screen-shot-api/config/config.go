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
	DBConnectionString string `json:"DB_CONNECTION_STRING" envconfig:"DB_CONNECTION_STRING"`
	DBEngine           string `json:"DB_ENGINE" envconfig:"DB_ENGINE"`
	WebAddress         string `json:"SCREEN_SHOT_API_ADDRESS" envconfig:"SCREEN_SHOT_API_ADDRESS"`
	WebPort            int    `json:"SCREEN_SHOT_API_PORT" envconfig:"SCREEN_SHOT_API_PORT"`
	ScreenShotServer   string `json:"SCREEN_SHOT_SERVER" envconfig:"SCREEN_SHOT_SERVER"`
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
