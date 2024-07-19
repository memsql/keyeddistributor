package keyed

import (
	"github.com/sharnoff/eventdistributor"

	"singlestore.com/helios/util/refcountmap"
)

/*
Package keyeddistributor provides a way to wrap [eventdistributor]
so that it is efficient to receive notification only for the
specific events you're interested in: those that have a specific
value as derived from the underlying event.

[eventdistributor]: https://pkg.go.dev/github.com/sharnoff/eventdistributor
*/

var closedChannel = func() <-chan struct{} {
	c := make(chan struct{})
	close(c)
	return c
}()

// Distributor supports subscribing to get a Reader
type Distributor[E any, K comparable] struct {
	m   *refcountmap.Map[K, *eventdistributor.Distributor[E]]
	key func(E) K
}

// Reader is used by a single subscriber to get events. Select on
// reader.WaitChan() to know when there is an event ready to consume
// and then use reader.Consume() to get the event. Use reader.Unsubscribe()
// when the reader is no longer needed. Reader embeds
// eventdistributor's Reader.
type Reader[E any] struct {
	eventdistributor.Reader[E]
	release func()
}

func New[E any, K comparable](key func(E) K) *Distributor[E, K] {
	return &Distributor[E, K]{
		m: refcountmap.New[K](func() *eventdistributor.Distributor[E] {
			return eventdistributor.New[E]()
		}),
		key: key,
	}
}

// Submit pushes an event into the Distributor. The returned channel is closed
// when the event has been fully consumed. If there are no subscribers, the
// event will be considered consumed immediately.
//
// Submit is thread-safe
func (k *Distributor[E, K]) Submit(e E) <-chan struct{} {
	key := k.key(e)
	dist, ok := k.m.Load(key)
	if !ok {
		return closedChannel
	}
	return dist.Submit(e)
}

// Subscribe creates a Reader that listens for events where the key matches
// a specific value. It is recommended that immediately after a Subscribe,
// that you defer the Unsubscribe:
//
//	reader := distributor.Subscribe(someValue)
//	defer reader.Unsubscribe()
//
// Subscribe is thread-safe
func (k *Distributor[E, K]) Subscribe(value K) *Reader[E] {
	d, release, _ := k.m.Get(value)
	reader := d.Subscribe()
	return &Reader[E]{
		Reader:  reader,
		release: release,
	}
}

// Unsubscribe releases a Reader. After an Unsubscribe, the WaitChan() and
// Consume() methods should not be used. It is important to Unsubscribe()
// because otherwise the Distributor will keep buffering events that are
// meant for the reader.
func (r *Reader[E]) Unsubscribe() {
	r.release()
	r.Reader.Unsubscribe()
}
