package handlers_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"receiver/messaging"
	"receiver/service"
	"receiver/tests/mocks"
	"strings"
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
)

//!+test
//go test -v

func getMockKafka(t *testing.T, messages []string) (*gomock.Controller, *sync.WaitGroup) {
	wg := &sync.WaitGroup{}
	mockCtrl := gomock.NewController(t)
	mockKafak := mocks.NewMockSyncProducer(mockCtrl, wg)
	var partition int32 = 1
	var offset int64 = 1
	for _, m := range messages {
		var err error
		msg := &sarama.ProducerMessage{
			Topic: "",
			Value: sarama.ByteEncoder(m),
		}
		if strings.Contains(m, "error") {
			err = errors.New("Kafka error")
		}
		wg.Add(1)
		mockKafak.EXPECT().SendMessage(msg).Return(partition, offset, err).Times(1)
	}
	messaging.InitializeKafka(mockKafak)
	return mockCtrl, wg
}

func TestReceiveJSON(t *testing.T) {
	mockCtrl, kafkawg := getMockKafka(t, []string{`{"url":"https://www.google.com/"}`, `{"url":"https://www.google.com/"}`, `{"url":"https://www.twitter.com/"}`, `{"url":"https://www.facebook.com/"}`})
	defer mockCtrl.Finish()

	tt := []struct {
		name            string
		body            string
		expectedMassage string
		statusCode      int
	}{
		{
			name:            "Send",
			body:            `{"url":"https://www.google.com/"}`,
			statusCode:      http.StatusCreated,
			expectedMassage: `{"message":"Success"}`,
		},
		{
			name:            "Send",
			body:            `[{"url":"https://www.google.com/"},{"url":"https://www.twitter.com/"},{"url":"https://www.facebook.com/"}]`,
			statusCode:      http.StatusCreated,
			expectedMassage: `{"message":"Success"}`,
		},
		{
			name:            "No Body",
			expectedMassage: `{"message":"unexpected end of JSON input"}`,
			statusCode:      http.StatusBadRequest,
		},
		{
			name:            "Error Json",
			body:            `the is not json`,
			expectedMassage: `{"message":"invalid character 'h' in literal true (expecting 'r')"}`,
			statusCode:      http.StatusBadRequest,
		},
		{
			name:            "Data Error",
			body:            `{"error":"data"}`,
			statusCode:      http.StatusBadRequest,
			expectedMassage: `{"message":"json: cannot unmarshal object into Go value of type []model.MessageModel"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "http://::/json", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()

			h := service.NewRouter()
			h.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("could not read response: %v", err)
			}

			if res.StatusCode != tc.statusCode {
				t.Errorf("expected status %v; got %v", tc.statusCode, res.Status)
			}

			if string(bytes.TrimSpace(b)) != tc.expectedMassage {
				t.Errorf("expected %v; got %s", tc.expectedMassage, b)
			}
		})
	}
	kafkawg.Wait()
}

func TestReady(t *testing.T) {
	tt := []struct {
		name           string
		kafkaReadiness bool
		expecte        string
		statusCode     int
	}{
		{name: "Ready", kafkaReadiness: true, expecte: `{"message":"Ready"}`, statusCode: http.StatusOK},
		{name: "Unready", kafkaReadiness: false, expecte: `{"message":"Unready"}`, statusCode: http.StatusServiceUnavailable},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			messaging.KafkaReadiness = tc.kafkaReadiness
			req, err := http.NewRequest("GET", "http://::/ready", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()

			h := service.NewRouter()
			h.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("could not read response: %v", err)
			}

			if res.StatusCode != tc.statusCode {
				t.Errorf("expected status %v; got %v", tc.statusCode, res.Status)
			}

			if string(bytes.TrimSpace(b)) != tc.expecte {
				t.Errorf("expected %v; got %s", tc.expecte, b)
			}
		})
	}
}

//!-tests

/*


//!+bench
//go test -v -bench=.

// TBR

*/
//!-bench
