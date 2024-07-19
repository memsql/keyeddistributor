package keyed_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"singlestore.com/helios/events/keyed"
)

type MyEvent struct {
	id int
}

func TestKeyed(t *testing.T) {
	distributor := keyed.New(func(e MyEvent) int {
		return e.id % 4
	})

	t.Log("subscribe to various keys")
	r0 := distributor.Subscribe(0)
	r1 := distributor.Subscribe(1)
	r5 := distributor.Subscribe(5) // should get no events

	t.Log("none are ready")
	require.False(t, ready(r0.WaitChan()))
	require.False(t, ready(r1.WaitChan()))
	require.False(t, ready(r5.WaitChan()))

	t.Log("submit an event that doesn't match any keys")
	s0 := distributor.Submit(MyEvent{id: 2}) // nobody
	require.False(t, ready(r0.WaitChan()))
	require.False(t, ready(r1.WaitChan()))
	require.False(t, ready(r5.WaitChan()))
	require.True(t, ready(s0))

	t.Log("submit an event that matches one key (0)")
	s1 := distributor.Submit(MyEvent{id: 0})
	require.True(t, ready(r0.WaitChan()))
	require.False(t, ready(r1.WaitChan()))
	require.False(t, ready(r5.WaitChan()))
	require.False(t, ready(s1))
	assert.Equal(t, MyEvent{id: 0}, r0.Consume())
	require.False(t, ready(r0.WaitChan()))
	require.False(t, ready(r1.WaitChan()))
	require.False(t, ready(r5.WaitChan()))
	require.True(t, ready(s1))

	t.Log("submit an event that matches one key (1)")
	s2 := distributor.Submit(MyEvent{id: 1})
	require.False(t, ready(r0.WaitChan()))
	require.True(t, ready(r1.WaitChan()))
	require.False(t, ready(r5.WaitChan()))
	require.False(t, ready(s2))
	assert.Equal(t, MyEvent{id: 1}, r1.Consume())
	require.False(t, ready(r0.WaitChan()))
	require.False(t, ready(r1.WaitChan()))
	require.False(t, ready(r5.WaitChan()))
	require.True(t, ready(s2))

	t.Log("note that unsubscribe releases the submission")
	s3 := distributor.Submit(MyEvent{id: 1})
	require.False(t, ready(r0.WaitChan()))
	require.True(t, ready(r1.WaitChan()))
	require.False(t, ready(r5.WaitChan()))
	require.False(t, ready(s3))
	r1.Unsubscribe()
	require.False(t, ready(r0.WaitChan()))
	require.False(t, ready(r5.WaitChan()))
	require.True(t, ready(s3))
}

func ready(c <-chan struct{}) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}
