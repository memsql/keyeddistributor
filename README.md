# keyeddistributor - distribute events to listeners by key

[![GoDoc](https://godoc.org/github.com/memsql/keyeddistributor?status.svg)](https://pkg.go.dev/github.com/memsql/keyeddistributor)
![unit tests](https://github.com/memsql/keyeddistributor/actions/workflows/go.yml/badge.svg)
[![report card](https://goreportcard.com/badge/github.com/memsql/keyeddistributor)](https://goreportcard.com/report/github.com/memsql/keyeddistributor)
[![codecov](https://codecov.io/gh/memsql/keyeddistributor/branch/main/graph/badge.svg)](https://codecov.io/gh/memsql/keyeddistributor)

Install:

	go get github.com/memsql/keyeddistributor

---

Keyeddistributor is an add-on to [eventdistributor](https://github.com/sharnoff/eventdistributor) that
allows for ad-hoc subscription to events that match a key value.

## An example

```go
keyed := keyeddistributor.New(func(t Thing) ThingKey) {
	return t.Key // or whatever is needed to extract a ThingKey from a Thing
})

reader := keyed.Subscribe("some key value")
defer reader.Unsubscribe()

<-reader.WaitChan()
thing := reader.Consume()
// do stuff with thing
```
