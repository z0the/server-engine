package eventbus

import (
	"sync"

	"go.uber.org/zap"
)

func NewBus(lg *zap.SugaredLogger) Bus {
	return &bus{
		lg:          lg,
		topics:      make(map[string]struct{}),
		subscribers: make(map[string][]MsgHandler),
	}
}

// There is no garbage collection, and no functionality for unsubscribe from topics
type bus struct {
	sync.RWMutex
	lg          *zap.SugaredLogger
	topics      map[string]struct{}
	subscribers map[string][]MsgHandler
}

func (b *bus) PublishWithToPrefix(topic string, msg BusMsg) {
	b.Publish(addToPrefix(topic), msg)
}

func (b *bus) PublishWithFromPrefix(topic string, msg BusMsg) {
	b.Publish(addFromPrefix(topic), msg)
}

func (b *bus) Publish(topic string, msg BusMsg) {
	if !b.isTopicExists(topic) {
		b.Lock()
		b.topics[topic] = struct{}{}
		b.Unlock()
	}

	b.RLock()
	defer b.RUnlock()
	for _, handler := range b.subscribers[topic] {
		go handler(msg)
	}
}

func (b *bus) SubscribeWithToPrefix(topic string, newSub MsgHandler) {
	b.Subscribe(addToPrefix(topic), newSub)
}

func (b *bus) SubscribeWithFromPrefix(topic string, newSub MsgHandler) {
	b.Subscribe(addFromPrefix(topic), newSub)
}

func (b *bus) Subscribe(topic string, newSub MsgHandler) {
	if !b.isTopicExists(topic) {
		b.Lock()
		b.topics[topic] = struct{}{}
		b.Unlock()
	}

	b.Lock()
	defer b.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], newSub)
}

func (b *bus) isTopicExists(subjectUID string) bool {
	b.RLock()
	defer b.RUnlock()
	_, isExists := b.topics[subjectUID]
	return isExists
}

func addToPrefix(topic string) string {
	return "to:" + topic
}

func addFromPrefix(topic string) string {
	return "from:" + topic
}
