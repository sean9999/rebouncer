package rebouncer

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
	Interrupt()
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
		incomingEvents: make(chan NICE, bufferSize),
		outgoingEvents: make(chan NICE, bufferSize),
		lifeCycle:      make(chan lifeCycleState, 1),
	}

	m.SetLifeCycleState(Running)

	//	ingest loop
	go func() {
		for niceEvent := range m.incomingEvents {
			m.reduce(reduceFunc, niceEvent)
		}
	}()

	//	quantize loop
	go func() {
		for lifeEvent := range m.lifeCycle {
			switch lifeEvent {
			case Quantizing:
				m.quantize(quantizeFunc)
			case Emiting:
				m.emit(emitFunc)
			case Draining:
				//close(m.incomingEvents)
				m.SetLifeCycleState(Draining)
			case Drained:
				close(m.outgoingEvents)
			}
		}
	}()

	m.lifeCycle <- Quantizing // start quantizer
	m.ingest(ingestFunc)

	return &m

}
