package parsing

import (
	"bytes"
	"encoding/json"
	"io"
	"receiver/logger"
	"receiver/model"
)

// JSON : kafka handler for JSON data
type jsonHandler struct {
	Topic string
}

// Handler :
func (J *jsonHandler) Handler(r io.Reader) error {
	data, err := unmarshalJSON(r)
	if err != nil {
		logger.Error(err)
		return err
	}
	// start send to kafka
	for _, v := range data {
		go sendDataToKafka(v, J.Topic)
	}
	// end send to kafka
	return nil
}

func unmarshalJSON(r io.Reader) ([]model.MessageModel, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	// check if the input is single objects.
	dataInputjson := model.MessageModel{}
	err := json.Unmarshal(buf.Bytes(), &dataInputjson)
	if err == nil && dataInputjson.URL != "" {
		return []model.MessageModel{dataInputjson}, nil
	}

	// check if the input is array of objects.
	arrdataInputjson := []model.MessageModel{}
	err = json.Unmarshal(buf.Bytes(), &arrdataInputjson)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return arrdataInputjson, nil
}
