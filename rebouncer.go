package rebouncer

import "fmt"

// all channels have this capacity
const DefaultBufferSize = 1024

// Rebouncer implements Behaviour
type Behaviour[NAUGHTY any, NICE any, BEAUTIFUL any] interface {
	Subscribe() <-chan BEAUTIFUL        // returns a channel and pushes events to it
	emit(EmitFunction[NICE, BEAUTIFUL]) // flushes the Queue
	readQueue() []NICE                  // gets the Queue, with safety and locking
	writeQueue([]NICE)                  // sets the Queue, handling safety and locking
	ingest(IngestFunction[NAUGHTY, NICE])
	quantize(QuantizeFunction[NICE])   // decides whether the flush the Queue
	reduce(ReduceFunction[NICE], NICE) // removes unwanted NiceEvents from the Queue
}

func NewRebouncer[NAUGHTY any, NICE any, BEAUTIFUL any](
	ingestFunc IngestFunction[NAUGHTY, NICE],
	reduceFunc ReduceFunction[NICE],
	quantizeFunc QuantizeFunction[NICE],
	emitFunc EmitFunction[NICE, BEAUTIFUL],
	bufferSize int,
) Behaviour[NAUGHTY, NICE, BEAUTIFUL] {

	//	channels
	m := stateMachine[NAUGHTY, NICE, BEAUTIFUL]{
		readyChannel:   make(chan bool),
		doneChannel:    make(chan bool),
		incomingEvents: make(chan NICE, bufferSize),
		outgoingEvents: make(chan BEAUTIFUL, bufferSize),
	}

	//	core functionality
	go func() {
		for {
			select {
			case niceEvent := <-m.incomingEvents:
				//m.Lock()
				m.reduce(reduceFunc, niceEvent) // for every event pushed to the queue, run reducer
				//m.Unlock()
			case isReady := <-m.readyChannel:
				//m.Lock()
				if isReady {
					m.emit(emitFunc)
				} else {
					m.quantize(quantizeFunc)
				}
				//m.Unlock()
			case doneVal := <-m.doneChannel:
				fmt.Println("all done", doneVal)
				m.incomingEvents = nil
				close(m.outgoingEvents)
			}
		}
	}()
	go m.ingest(ingestFunc) // start ingesting
	m.readyChannel <- false // start quantizer

	return &m

}
