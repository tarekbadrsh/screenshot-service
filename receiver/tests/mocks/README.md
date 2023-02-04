
# To Generate 'SyncProducer' run this command
`mockgen -destination=mocks/mock_SyncProducer.go -package=mocks github.com/Shopify/sarama SyncProducer`

or you can put this comment above `var kafkaProducer sarama.SyncProducer` in `messaging/producer.go` file and run `go generate` inside `messaging` directory 
```
//go:generate mockgen -destination=../mocks/mock_SyncProducer.go -package=mocks github.com/Shopify/sarama SyncProducer
var kafkaProducer sarama.SyncProducer
```

## Please Note: the generated code should be updated with 'WaitGroup' object to be able to handle the concurrency submit to kafka

e.g.
### update `MockSyncProducer`
```
type MockSyncProducer struct {
    ...
	wg       *sync.WaitGroup
}
```

### update `NewMockSyncProducer`
```
func NewMockSyncProducer(ctrl *gomock.Controller, wg *sync.WaitGroup) *MockSyncProducer {
	mock := &MockSyncProducer{ctrl: ctrl, wg: wg}
    ...
}
```

### update `func SendMessage`
```
func (m *MockSyncProducer) SendMessage(arg0 *sarama.ProducerMessage) (int32, int64, error) {
    ...
	m.wg.Done()
	return ret0, ret1, ret2
}
```