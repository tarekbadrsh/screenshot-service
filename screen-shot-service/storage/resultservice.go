package storage

import (
	"bytes"
	"encoding/json"
	"net/http"
	"screen-shot-service/logger"
	"screen-shot-service/model"

	"screen-shot-service/config"
	"sync"
)

var (
	resultServiceURL string
	initResSrvOnce   sync.Once
)

// InitializeResultService : initialize the communication with the result service
func InitializeResultService(c config.Config) {
	initResSrvOnce.Do(func() {
		resultServiceURL = c.ResultServiceURL
		logger.Info("Result service has been initialized")
	})
}

// SendGeneratedResult : send the result of generate the screen-shot to results service
func SendGeneratedResult(result *model.GeneratorResult) error {
	jsonValue, _ := json.Marshal(result)
	_, err := http.Post(resultServiceURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
