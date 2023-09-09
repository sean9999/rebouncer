# Rebouncer

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/rebouncer/graphs/commit-activity)

[![Go Reference](https://pkg.go.dev/badge/github.com/sean9999/rebouncer.svg)](https://pkg.go.dev/github.com/sean9999/rebouncer)

[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/rebouncer)](https://goreportcard.com/report/github.com/sean9999/rebouncer)

[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/rebouncer.svg)](https://github.com/sean9999/rebouncer)

## A Powerful Debouncer for your Conjuring Needs

Rebouncer is a generic library that takes a noisy source of events, and produces a cleaner source. It does debouncing, and much more, by offering a well-defined set of lifecycle states and ways to hook into them.

## Concepts

Rebouncer has two run-loops:

### Ingest => Reduce

The Ingestor runs in it's own loop, pushing events to a channel in Rebouncer. Every time an event is pushed, Reducer runs. Reducer operates on the entire queue of events, filtering out unwanted events or modifying to taste. Here are the definitions of these functions. NICE is a type parameter. Internally, your custom event type is known as a "Nice Event".

```go
type Ingester[NICE any] func(chan<- NICE)
type Reducer[NICE any] func([]NICE) []NICE
```

### Quantize => Emit

Quantizer returns true or false. True when we want to flush the queue to the consumer, and False when we don't. As soon as Quantizer is returned, it's run again. So to throttle it's behaviour, do `time.Sleep()`.

When the program enters the Draining state, it shuts down after the last Emit(). Otherwise it keeps looping.

```go
type Quantizer[NICE any] func([]NICE) bool
```

Rebouncer is generic. The atomic unit "event" is whatever shape you need it to be. Just make sure that your ingestor, reducer, and quantizer all operate on the same type.

```go
//	Example

type myEvent struct {
	id int
	name string
	timestamp time.Time
}

//	ingest events
ingest := func(incoming<- myEvent) {
	for ev := range mySourceOfEvents() {
		incoming<-ev
	} 
}

//	we're not interested in any event involving .DS_Store
reduce := func(inEvents []myEvent) []myEvent {
	outEvents := []myEvent{}
	for ev := range inEvents {
		if ev.name != ".DS_Store" {
			outEvents = append(outEvents, ev)
		}
	}
	return outEvents
}

//	flush the queue every second
quantize := func(queue []myEvent) bool {
	time.Sleep(time.Second)
	if len(queue) > 0 {
		return true
	} else {
		return false
	}
}

re := rebouncer.NewRebouncer[myEvent](ingest, reduce, quantize, 1024)

for ev := range re.Subscribe() {
	fmt.Println(ev)
}
```

<img src="/docs/hand.jpg" width="450" />
