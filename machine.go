package rebouncer

import (
	"fmt"
	"sync"
)

// machine is the data structure that allows *machine to implement [Behaviour]
type machine[T any] struct {
	outgoingEvents chan NiceEvent[T] // NiceEvents for our consumer
	readyChannel   chan bool
	Queue          []NiceEvent[T] // an intermediary storage mechanism
	bufferSize     int            // all channels should have this buffer size
	mu             sync.Mutex     // lock batchMap when we're processing it
	ingestor       IngestFunction[T]
	reducer        ReduceFunction[T]
	quantizer      QuantizeFunction[T]
}

// Config contains everything you need to pass to the New() constructor
type Config[T any] struct {
	BufferSize int
	Ingestor   IngestFunction[T]
	Reducer    ReduceFunction[T]
	Quantizer  QuantizeFunction[T]
}

// New is the canonical way to construct a new StateMachine
func New[T any](config Config[T]) (Behaviour[T], error) {

	emptyQueue := make([]NiceEvent[T], 0, config.BufferSize)

	m := machine[T]{
		outgoingEvents: make(chan NiceEvent[T], config.BufferSize),
		readyChannel:   make(chan bool, config.BufferSize),
		bufferSize:     config.BufferSize,
		Queue:          emptyQueue,
		quantizer:      config.Quantizer,
		reducer:        config.Reducer,
		ingestor:       config.Ingestor,
	}

	//	incomingEvents is a channel that is being pushed to by [IngestFunction]
	incomingEvents, err := m.ingest()
	if err != nil {
		return nil, err
	}

	//	readyChan is a channel that [QuantizeFunction] writes true to when it's time to [emit]
	//	emit() will transfer NiceEvents from Queue to outgoingEvents
	//m.readyChannel = m.quantize()

	//	listen to events emitted by Ingestor
	go func() {
		for inEvent := range incomingEvents {
			m.reduce(inEvent)
		}
	}()

	//	send one value to readyChan to get the ball rolling
	go m.quantize()

	//	keep the ball rolling
	go func() {
		for isReady := range m.readyChannel {
			fmt.Printf("isReady is %v\n", isReady)
			if isReady {
				m.emit()
			} else {
				m.quantizer(m.readyChannel, m.readQueue())
			}
		}
	}()

	return &m, nil
}

func (m *machine[T]) readQueue() []NiceEvent[T] {
	return m.Queue
}

func (m *machine[T]) writeQueue(eventSlice []NiceEvent[T]) error {
	m.mu.Lock()
	m.Queue = eventSlice
	m.mu.Unlock()
	return nil
}

// Emit all queued NiceEvents to OutgoingEvents
func (m *machine[T]) emit() {
	//m.mu.Lock()

	fmt.Println("emit()")

	for _, e := range m.Queue {
		m.outgoingEvents <- e
	}
	m.Queue = []NiceEvent[T]{}
	//m.mu.Unlock()
}

// Subscribe gives us our final, curated channel of NiceEvents
func (m *machine[T]) Subscribe() <-chan NiceEvent[T] {
	return m.outgoingEvents
}
