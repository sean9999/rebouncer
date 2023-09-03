package rebouncer

import "fmt"

// all channels have this capacity
const DefaultBufferSize = 1024

type Rebouncer[NICE any] interface {
	Subscribe() <-chan NICE  // returns a channel and pushes events to it
	emit(EmitFunction[NICE]) // flushes the Queue
	readQueue() []NICE       // gets the Queue, with safety and locking
	writeQueue([]NICE)       // sets the Queue, handling safety and locking
	ingest(IngestFunction[NICE])
	quantize(QuantizeFunction[NICE])   // decides whether the flush the Queue
	reduce(ReduceFunction[NICE], NICE) // removes unwanted NiceEvents from the Queue
	Drain()
}

func NewRebouncer[NICE any](
	ingestFunc IngestFunction[NICE],
	reduceFunc ReduceFunction[NICE],
	quantizeFunc QuantizeFunction[NICE],
	emitFunc EmitFunction[NICE],
	bufferSize int,
) Rebouncer[NICE] {

	//	channels
	m := stateMachine[NICE]{
		readyChannel:   make(chan bool), //	indicates we're ready to flush to client
		doneChannel:    make(chan bool), //	indicates we're done consuming events and ready to drain
		incomingEvents: make(chan NICE, bufferSize),
		outgoingEvents: make(chan NICE, bufferSize),
	}

	m.lifeCycleState = Running

	//	core functionality
	go func() {
		for m.lifeCycleState != Drained {
			select {
			case niceEvent := <-m.incomingEvents:
				m.reduce(reduceFunc, niceEvent) // for every event pushed to the queue, run reducer
			case isReady := <-m.readyChannel:
				if isReady {
					m.emit(emitFunc)
				} else {
					m.quantize(quantizeFunc)
				}
			case doneVal := <-m.doneChannel:
				fmt.Println("all done just kidding", doneVal)
				m.Drain()
				/*
					for _, niceEvent := range m.queue {
						m.outgoingEvents <- emitFunc(niceEvent)
					}
				*/
			}
		}
	}()
	m.readyChannel <- false // start quantizer
	m.ingest(ingestFunc)    // start ingesting

	return &m

}
