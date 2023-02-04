package parsing

import (
	"encoding/json"
	"io"
	"receiver/logger"
	"receiver/messaging"
	"receiver/model"
)

const (
	// JSON for json data
	JSON = iota

	// add more data format here
)

// KafkaHandler :
type KafkaHandler interface {
	Handler(r io.Reader) error
}

// GetHandler :
func GetHandler(handler int, topic string) KafkaHandler {
	switch handler {
	case JSON: // to handle json data
		return &jsonHandler{Topic: topic}
		// we could add more handler to handle any data format
	}
	return nil
}

func sendDataToKafka(msg model.MessageModel, topic string) error {
	b, err := json.Marshal(msg)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = messaging.Produce(topic, b)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
