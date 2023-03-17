package rebouncer

import (
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
}

// pointer to machinery implements StateMachine
type machinery struct {
	OutgoingEvents chan NiceEvent
	batchMap       map[string]NiceEvent
	bufferSize     int
	mu             sync.Mutex
}

// pass this in to the New() constructor
type Config struct {
	BufferSize int
}

// The easiest way to create a new StateMachine
func New(config Config) StateMachine {
	m := machinery{
		OutgoingEvents: make(chan NiceEvent),
		bufferSize:     config.BufferSize,
		batchMap:       map[string]NiceEvent{},
	}
	return &m
}

// Injest takes a NiceEvent and either appends it to batchMap or ignores it
//
//	Additionally, it decides whether to call Emit() or not
func (m *machinery) Injest(newEvent NiceEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

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

	if len(m.batchMap) > 3 {
		m.Emit()
	}

}

// Emits all the queued NiceEvents to OutgoingEvents
func (m *machinery) Emit() {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	for _, niceEvent := range m.batchMap {
		m.OutgoingEvents <- niceEvent
	}
	m.batchMap = map[string]NiceEvent{}
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
