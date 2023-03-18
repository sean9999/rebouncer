package rebouncer

import (
	"sync"
)

// StateMachine respresents access all necessary methods. It is the entrypoint for the app.
// The [New] constructor should be used to create instances which can fulfil these methods.
type StateMachine interface {
	Subscribe() chan NiceEvent
	Version() string
	Info() map[string]any
	WatchDir(string)
	Emit()
	GetQueue() []NiceEvent
	SetQueue([]NiceEvent) error
	Push(NiceEvent)
	Quantize(chan bool, *[]NiceEvent)
	Reduce([]NiceEvent) []NiceEvent
}

type userFunctions struct {
	quantizer Quantizer
	reducer   Reducer
	injestor  Injestor
}

// pointer to machinery implements StateMachine
type machinery struct {
	OutgoingEvents chan NiceEvent // NiceEvents for our consumer
	readyChan      chan bool      // whe true is sent here, a batch is ready
	Queue          []NiceEvent    // an intermediary storage mechanism
	bufferSize     int            // all channels should have this buffer size
	mu             sync.Mutex     // lock batchMap when we're processing it
	user           userFunctions  // user passes in these functions
}

// pass this in to the New() constructor
type Config struct {
	BufferSize int
	Quantizer  Quantizer
	Reducer    Reducer
	Injestor   Injestor
}

// The canonical way to create a new StateMachine
func New(config Config) StateMachine {

	m := machinery{
		OutgoingEvents: make(chan NiceEvent, config.BufferSize),
		readyChan:      make(chan bool, config.BufferSize),
		bufferSize:     config.BufferSize,
		Queue:          []NiceEvent{},
		user: userFunctions{
			quantizer: config.Quantizer,
			reducer:   config.Reducer,
			injestor:  config.Injestor,
		},
	}

	incomingEvents := m.user.injestor()

	//	listen to events emitted by Injestor
	go func() {
		for inEvent := range incomingEvents {
			m.Push(inEvent)
		}
	}()

	//	Emit() whenever we get a 'true' value on readyChan
	go func() {
		for isReady := range m.readyChan {
			if isReady {
				m.Emit()
			}
		}
	}()

	return &m
}

func (m *machinery) GetQueue() []NiceEvent {
	return m.Queue
}

func (m *machinery) SetQueue(eventSlice []NiceEvent) error {
	m.mu.Lock()
	m.Queue = eventSlice
	m.mu.Unlock()
	return nil
}

func (m *machinery) Reduce(inEvents []NiceEvent) []NiceEvent {
	outEvents := m.user.reducer(inEvents)
	return outEvents
}

// Push takes a NiceEvent and either appends it to batchMap or ignores it
//
//	Additionally, it decides whether to call Emit() or not
func (m *machinery) Push(newEvent NiceEvent) {
	m.Queue = append(m.Queue, newEvent)
	m.Queue = m.Reduce(m.Queue)
	go m.Quantize(m.readyChan, &m.Queue)
}

// Quantize runs after Injest() and decides whether or not to call Emit()
func (m *machinery) Quantize(readyChannel chan bool, em *[]NiceEvent) {
	fn := m.user.quantizer
	go fn(readyChannel, em)
}

// Emits all the queued NiceEvents to OutgoingEvents
func (m *machinery) Emit() {
	//m.mu.Lock()
	for _, e := range m.Queue {
		m.OutgoingEvents <- e
	}
	m.Queue = []NiceEvent{}
	//m.mu.Unlock()
}

// the channel of NiceEvents we want to operate on
func (m *machinery) Subscribe() chan NiceEvent {
	return m.OutgoingEvents
}

// returns the application version string
func (m *machinery) Version() string {
	return appVersion
}

// returns useful (maybe?) information
func (m *machinery) Info() map[string]any {
	r := map[string]any{
		"bufferSize": m.bufferSize,
		"version":    m.Version(),
	}
	return r
}
