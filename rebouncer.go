package rebouncer

// all channels have this capacity
const DefaultBufferSize = 1024

// Rebouncer implements Behaviour
type Behaviour[NAUGHTY any, NICE any, BEAUTIFUL any] interface {
	Subscribe() <-chan BEAUTIFUL        // returns a channel and pushes events to it
	emit(EmitFunction[NICE, BEAUTIFUL]) // flushes the Queue
	readQueue() []NICE                  // gets the Queue, with safety and locking
	writeQueue([]NICE)                  // sets the Queue, handling safety and locking
	ingest(IngestFunction[NAUGHTY, NICE])
	stopIngesting() error
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

	m := stateMachine[NAUGHTY, NICE, BEAUTIFUL]{
		readyChannel:   make(chan bool),
		incomingEvents: make(chan NICE, bufferSize),
		outgoingEvents: make(chan BEAUTIFUL, bufferSize),
	}

	//	core functionality
	go func() {
		m.ingest(ingestFunc)
		m.readyChannel <- false
		for {
			select {
			case niceEvent := <-m.incomingEvents:
				m.reduce(reduceFunc, niceEvent)
			case isReady := <-m.readyChannel:
				if isReady {
					m.emit(emitFunc)
					m.readyChannel <- false
				} else {
					m.quantize(quantizeFunc)
				}
			}
		}
	}()

	return &m

}
