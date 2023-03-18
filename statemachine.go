package rebouncer

import (
	"fmt"
	"sync"
	"time"
)

// StateMachine is used to access all necessary methods and data
type StateMachine interface {
	Subscribe() chan NiceEvent
	Version() string
	Info() map[string]any
	WatchDir(string)
	Emit()
	Injest(NiceEvent)
	Quantize(chan bool, *EventMap)
}

type userFunctions struct {
	quantizer Quantizer
}

// pointer to machinery implements StateMachine
type machinery struct {
	OutgoingEvents chan NiceEvent // NiceEvents for our consumer
	readyChan      chan bool      // whe true is sent here, a batch is ready
	batchMap       EventMap       // an intermediary storage mechanism
	bufferSize     int            // all channels should have this buffer size
	mu             sync.Mutex     // lock batchMap when we're processing it
	userFuncs      userFunctions  // user passes in these functions
	ticker         time.Ticker
}

// pass this in to the New() constructor
type Config struct {
	BufferSize int
	Quantizer  Quantizer
}

// The easiest way to create a new StateMachine
func New(config Config) StateMachine {

	m := machinery{
		OutgoingEvents: make(chan NiceEvent, 1024),
		readyChan:      make(chan bool, 1024),
		bufferSize:     config.BufferSize,
		batchMap:       EventMap{},
		userFuncs: userFunctions{
			quantizer: config.Quantizer,
		},
		ticker: *time.NewTicker(5 * time.Minute),
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

// Injest takes a NiceEvent and either appends it to batchMap or ignores it
//
//	Additionally, it decides whether to call Emit() or not
func (m *machinery) Injest(newEvent NiceEvent) {
	//m.mu.Lock()
	//defer m.mu.Unlock()

	fmt.Println("injest")

	fileName := newEvent.File
	existingEvent, existsInMap := m.batchMap[fileName]

	//	only process non-temp files and valid events
	if !isTempFile(fileName) && !newEvent.IsZeroed() {
		if existsInMap {
			switch {
			case existingEvent.Operation == "notify.Create" && newEvent.Operation == "notify.Remove":
				//	a Create followed by a Remove means nothing of significance happened. Purge the record
				delete(m.batchMap, fileName)
			default:
				//	the default case should be to overwrite the record
				newEvent.Topic = "rebouncer/outgoing/1"
				m.batchMap[fileName] = newEvent
			}
		} else {
			newEvent.Topic = "rebouncer/outgoing/0"
			m.batchMap[newEvent.File] = newEvent
		}
	}

	go m.Quantize(m.readyChan, &m.batchMap)

}

// Quantize runs after Injest() and decides whether or not to call Emit()
func (m *machinery) Quantize(readyChannel chan bool, em *EventMap) {
	fmt.Println("Quantize()")
	fn := m.userFuncs.quantizer
	go fn(readyChannel, em)
}

// Emits all the queued NiceEvents to OutgoingEvents
func (m *machinery) Emit() {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	for _, niceEvent := range m.batchMap {
		m.OutgoingEvents <- niceEvent
	}
	fmt.Println("Emit()")
	m.batchMap = EventMap{}
	fmt.Println(m.batchMap)
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
