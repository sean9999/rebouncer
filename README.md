# Rebouncer

## A Powerful Debouncer for your Conjuring Needs

<div class="evenly-spread">

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/rebouncer/graphs/commit-activity) 
[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/rebouncer)](https://goreportcard.com/report/github.com/sean9999/rebouncer)
[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/rebouncer.svg)](https://github.com/sean9999/rebouncer)
[![Go Reference](https://pkg.go.dev/badge/github.com/sean9999/rebouncer.svg)](https://pkg.go.dev/github.com/sean9999/rebouncer)

<div>

<img src="/docs/hand.jpg" width="350" />

Rebouncer is a package that takes a noisy source of events and produces a cleaner, fitter, happier (not drinking too much) source. It is useful in scenarios where you want debounce-like functionality, and full control over how events are consumed, filtered, queued and flushed to the consuming process.

## Concepts

### The NICE Event

NiceEvent is simply the type of event you pass in to Rebouncer. It exists only as a concept so we can have something to refer to:

```go
type Rebouncer[NICE any] interface {
	Subscribe() <-chan NICE 	// the channel a consumer can subsribe to
	emit()                  	// flushes the Queue
	readQueue() []NICE      	// gets the Queue, with safety and locking
	writeQueue([]NICE)      	// sets the Queue, handling safety and locking
	ingest(Ingester[NICE])
	quantize(Quantizer[NICE])   // decides whether the flush the Queue
	reduce(Reducer[NICE], NICE) // removes unwanted NiceEvents from the Queue
	Interrupt()                 //	call this to initiate the "Draining" state
}

type myType struct {
	...
}

bufferSize = 1024 // how much buffer space do we want for incoming events?

//	myRebouncer is a Rebouncer of type myType
myRebouncer := rebouncer.NewRebouncer[myType](ingest, reduce, quantize, bufferSize)
```

Rebouncer has two run-loops:

### Ingest ☞ Reduce

The Ingestor runs in it's own loop, pushing events to a channel in Rebouncer. Every time an event is pushed, Reducer runs. Reducer operates on the entire queue of events, filtering out unwanted events or modifying to taste. Here are the definitions of these functions. NICE is a type parameter. Internally, your custom event type is known as a "Nice Event".

```go
type Ingester[NICE any] func(chan<- NICE)
type Reducer[NICE any] func([]NICE) []NICE
```

### Quantize ☞ Emit

Quantizer returns true or false. True when we want to flush the queue to the consumer, and False when we don't. As soon as Quantizer is returned, it's run again. So to throttle it, do `time.Sleep()`.

When the program enters the Draining state, it shuts down after the last Emit(). Otherwise it keeps looping.

```go
type Quantizer[NICE any] func([]NICE) bool
```

Ensure that your Ingestor, Reducer, and Quantizer all operate on the same type:

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


