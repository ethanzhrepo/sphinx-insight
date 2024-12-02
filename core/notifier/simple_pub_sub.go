package notifier

import "sync"

type SimplePubSub struct {
	sync.RWMutex
	subscribers map[string][]chan string
}

func NewSimplePubSub() *SimplePubSub {
	return &SimplePubSub{
		subscribers: make(map[string][]chan string),
	}
}

func (ps *SimplePubSub) Subscribe(topic string) chan string {
	ps.Lock()
	defer ps.Unlock()

	ch := make(chan string)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)

	return ch
}

func (ps *SimplePubSub) Unsubscribe(topic string, ch chan string) {
	ps.Lock()
	defer ps.Unlock()

	if chans, found := ps.subscribers[topic]; found {
		for i, sub := range chans {
			if sub == ch {
				ps.subscribers[topic] = append(chans[:i], chans[i+1:]...)
				close(sub)
				break
			}
		}
	}
}

func (ps *SimplePubSub) Publish(topic string, msg string) {
	ps.RLock()
	defer ps.RUnlock()

	if chans, found := ps.subscribers[topic]; found {
		for _, ch := range chans {
			go func(c chan string) {
				c <- msg
			}(ch)
		}
	}
}

func (ps *SimplePubSub) Close() {
	ps.Lock()
	defer ps.Unlock()

	for _, chans := range ps.subscribers {
		for _, ch := range chans {
			close(ch)
		}
	}
}
