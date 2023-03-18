package rebouncer

import (
	"fmt"
	"sync"
)

// StateMachine is used to access all necessary methods and data
type StateMachine interface {
	Subscribe() chan NiceEvent
	Version() string
	Info() map[string]any
	WatchDir(string)
	Emit()
	Injest(NiceEvent)
	Quantize(chan bool, *[]NiceEvent)
	Reduce([]NiceEvent) []NiceEvent
}

type userFunctions struct {
	quantizer Quantizer
	reducer   Reducer
}

// pointer to machinery implements StateMachine
type machinery struct {
	OutgoingEvents chan NiceEvent // NiceEvents for our consumer
	readyChan      chan bool      // whe true is sent here, a batch is ready
	batchArray     []NiceEvent    // an intermediary storage mechanism
	bufferSize     int            // all channels should have this buffer size
	mu             sync.Mutex     // lock batchMap when we're processing it
	user           userFunctions  // user passes in these functions
}

// pass this in to the New() constructor
type Config struct {
	BufferSize int
	Quantizer  Quantizer
	Reducer    Reducer
}

// The canonical way to create a new StateMachine
func New(config Config) StateMachine {

	m := machinery{
		OutgoingEvents: make(chan NiceEvent, config.BufferSize),
		readyChan:      make(chan bool, config.BufferSize),
		bufferSize:     config.BufferSize,
		batchArray:     []NiceEvent{},
		user: userFunctions{
			quantizer: config.Quantizer,
			reducer:   config.Reducer,
		},
	}

	//	Emit() whenever we get true on readyChan
	go func() {
		for isReady := range m.readyChan {
			if isReady {
				m.Emit()
			}
		}
	}()

	return &m
}

func (m *machinery) Shlock() {
	m.mu.Lock()
	fmt.Println("shlock")
	m.mu.Unlock()
}

func (m *machinery) Reduce(inEvents []NiceEvent) []NiceEvent {
	outEvents := m.user.reducer(inEvents)
	return outEvents
}

// Injest takes a NiceEvent and either appends it to batchMap or ignores it
//
//	Additionally, it decides whether to call Emit() or not
func (m *machinery) Injest(newEvent NiceEvent) {

	m.batchArray = append(m.batchArray, newEvent)
	m.batchArray = m.Reduce(m.batchArray)
	go m.Quantize(m.readyChan, &m.batchArray)

}

// Quantize runs after Injest() and decides whether or not to call Emit()
func (m *machinery) Quantize(readyChannel chan bool, em *[]NiceEvent) {
	fmt.Println("Quantize()")
	fn := m.user.quantizer
	go fn(readyChannel, em)
}

// Emits all the queued NiceEvents to OutgoingEvents
func (m *machinery) Emit() {
	for _, e := range m.batchArray {
		m.OutgoingEvents <- e
	}
	m.batchArray = []NiceEvent{}
}

func (m *machinery) Subscribe() chan NiceEvent {
	return m.OutgoingEvents
}
func (m *machinery) Version() string {
	return appVersion
}
func (m *machinery) Info() map[string]any {
	r := map[string]any{
		"bufferSize": m.bufferSize,
	}
	return r
}
