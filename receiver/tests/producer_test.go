package handlers_test

//!+test
//go test -v

// func TestKafkaProducer(t *testing.T) {
// 	tt := []struct {
// 		name      string
// 		topic     string
// 		message   string
// 		partition int32
// 		offset    int64
// 		err       error
// 	}{
// 		{name: "Message1", topic: "topic1", message: "message1", partition: 0, offset: 1, err: nil},
// 		{name: "Message2", topic: "topic2", message: "message2", partition: 1, offset: 2, err: nil},
// 		{name: "Message-error", topic: "topic3", message: "message3", partition: 2, offset: 3, err: errors.New("Kafka-error")},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			wg := &sync.WaitGroup{}
// 			logger.Info(tc.name)
// 			mockCtrl := gomock.NewController(t)
// 			defer mockCtrl.Finish()

// 			msg := &sarama.ProducerMessage{
// 				Topic: tc.topic,
// 				Value: sarama.ByteEncoder(tc.message),
// 			}
// 			mockKafak := mocks.NewMockSyncProducer(mockCtrl, wg)
// 			wg.Add(1)
// 			mockKafak.EXPECT().SendMessage(msg).Return(tc.partition, tc.offset, tc.err).Times(1)

// 			messaging.InitializeKafka(mockKafak)

// 			result := messaging.Produce(tc.topic, []byte(tc.message))
// 			if result != tc.err {
// 				t.Errorf("expected kakfa Error %v; got %v", result, tc.err)
// 			}
// 			wg.Wait()
// 		})
// 	}
// }

//!-tests
