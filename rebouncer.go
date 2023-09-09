package rebouncer

import (
	"sync"
)

// incomingEvents will have this capacity
const DefaultBufferSize = 1024

type Rebouncer[NICE any] interface {
	Subscribe() <-chan NICE // the channel a consumer can subsribe to
	emit()                  // flushes the Queue
	readQueue() []NICE      // gets the Queue, with safety and locking
	writeQueue([]NICE)      // sets the Queue, handling safety and locking
	ingest(Ingester[NICE])
	quantize(Quantizer[NICE])   // decides whether the flush the Queue
	reduce(Reducer[NICE], NICE) // removes unwanted NiceEvents from the Queue
	Interrupt()                 //	call this to initiate the "Draining" state
}

// NewRebouncer is the best way to create a new Rebouncer.
func NewRebouncer[NICE any](
	ingestFunc Ingester[NICE],
	reduceFunc Reducer[NICE],
	quantizeFunc Quantizer[NICE],
	bufferSize int, // for sizing the buffered channel that accepts incoming events
) Rebouncer[NICE] {

	//	channels
	m := rebounceMachine[NICE]{
		incomingEvents: make(chan NICE, bufferSize),
		outgoingEvents: make(chan NICE),
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
				m.emit()
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

// rebounceMachine implements [Rebouncer]
type rebounceMachine[NICE any] struct {
	lifeCycle      chan lifeCycleState
	incomingEvents chan NICE
	outgoingEvents chan NICE
	queue          Queue[NICE]
	mu             sync.RWMutex
	lifeState      lifeCycleState
}

func (m *rebounceMachine[NICE]) SetLifeCycleState(s lifeCycleState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lifeState = s
}

func (m *rebounceMachine[NICE]) GetLifeCycleState() lifeCycleState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lifeState
}

func (m *rebounceMachine[NICE]) Interrupt() {
	m.lifeCycle <- Draining
}

func (m *rebounceMachine[NICE]) Subscribe() <-chan NICE {
	return m.outgoingEvents
}
