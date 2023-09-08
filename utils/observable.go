package utils

import "sync"

type Subscription[V any] interface {
	Unsubscribe()
	Ch() <-chan V
}

type SubscriptionImpl[V any] struct {
	ch                   chan V
	closed               bool
	removeFromObservable func()
}

func (s *SubscriptionImpl[V]) Unsubscribe() {
	close(s.ch)
	s.removeFromObservable()
}

func (s *SubscriptionImpl[V]) Ch() <-chan V {
	return s.ch
}

type Observable[V any] interface {
	Subscribe() Subscription[V]
}

type ObservableImpl[V any] struct {
	mu          sync.RWMutex
	subscribers []chan V
	ch          <-chan V
	closed      bool
}

func NewControlledObservable[V any]() (Observable[V], chan V) {
	emitter := make(chan V, 1)
	o := &ObservableImpl[V]{sync.RWMutex{}, []chan V{}, emitter, false}
	o.listen(emitter)

	return o, emitter
}

func (o *ObservableImpl[V]) Subscribe() Subscription[V] {
	o.mu.Lock()
	defer o.mu.Unlock()

	ch := make(chan V, 1)
	o.subscribers = append(o.subscribers, ch)

	removeFromObservable := func() {
		o.mu.Lock()
		defer o.mu.Unlock()

		for i, ch := range o.subscribers {
			if ch == ch {
				o.subscribers = append(o.subscribers[:i], o.subscribers[i+1:]...)
				break
			}
		}
	}

	return &SubscriptionImpl[V]{ch, o.closed, removeFromObservable}
}

func (o *ObservableImpl[V]) listen(emitter <-chan V) {
	go func() {
		for v := range emitter {
			o.mu.RLock()
			for _, ch := range o.subscribers {
				ch <- v
			}
			o.mu.RUnlock()
		}
	}()
}
