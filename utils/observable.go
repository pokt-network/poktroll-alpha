package utils

import "sync"

// Observable is a generic interface that allows multiple subscribers to read from a single channel
type Observable[V any] interface {
	Subscribe() Subscription[V]
}

type ObservableImpl[V any] struct {
	mu          sync.RWMutex
	ch          <-chan V // private channel that is used to emit values to subscribers
	subscribers []chan V // subscribers is a list of channels that are subscribed to the observable
	closed      bool
}

// Creates a new observable which emissions are controlled by the emitter channel
func NewControlledObservable[V any](emitter chan V) (Observable[V], chan V) {
	// If the caller does not provide an emitter, create a new one and return it
	e := make(chan V, 1)
	if emitter != nil {
		e = emitter
	}
	o := &ObservableImpl[V]{sync.RWMutex{}, e, []chan V{}, false}

	// Start listening to the emitter and emit values to subscribers
	go o.listen(e)

	return o, e
}

// Get a subscription to the observable
func (o *ObservableImpl[V]) Subscribe() Subscription[V] {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Create a channel for the subscriber and append it to the subscribers list
	ch := make(chan V, 1)
	o.subscribers = append(o.subscribers, ch)

	// Removal function used when unsubscribing from the observable
	removeFromObservable := func() {
		o.mu.Lock()
		defer o.mu.Unlock()

		for i, s := range o.subscribers {
			if ch == s {
				o.subscribers = append(o.subscribers[:i], o.subscribers[i+1:]...)
				break
			}
		}
	}

	// Subscription gets its closed state from the observable
	return &SubscriptionImpl[V]{ch, o.closed, removeFromObservable}
}

// Listen to the emitter and emit values to subscribers
// This function is blocking and should be run in a goroutine
func (o *ObservableImpl[V]) listen(emitter <-chan V) {
	for v := range emitter {
		// Lock for o.subscribers slice as it can be modified by subscribers
		o.mu.RLock()
		for _, ch := range o.subscribers {
			ch <- v
		}
		o.mu.RUnlock()
	}

	// Here we know that the emitter has been closed, all subscribers should be closed as well
	o.mu.Lock()
	o.closed = true
	for _, ch := range o.subscribers {
		close(ch)
	}
	o.subscribers = []chan V{}
	o.mu.Unlock()
}

// Subscription is a generic interface that provide access to the underlying channel
// and allows unsubscribing from an observable
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
	s.closed = true
	s.removeFromObservable()
}

func (s *SubscriptionImpl[V]) Ch() <-chan V {
	return s.ch
}
