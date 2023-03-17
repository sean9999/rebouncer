package rebouncer

import "sync"

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
	NiceChannel chan NiceEvent
	batch       []NiceEvent
	bufferSize  int
	mu          sync.Mutex
}

// pass this in to the New() constructor
type Config struct {
	BufferSize int
}

// The easiest way to create a new StateMachine
func New(config Config) StateMachine {
	m := machinery{
		NiceChannel: make(chan NiceEvent),
		bufferSize:  config.BufferSize,
	}
	return &m
}

// Injest takes a NiceEvent and appends it to *batch
//
//	It may contain logic to filter out events we're not interested in
//
// It should also contain logic to decide whether to Emit() (to send along a bunch of NiceEvents to the consumer)
func (m *machinery) Injest(e NiceEvent) {
	m.batch = append(m.batch, e)

	//	let's say when there are 15 events, we Emit()

	if len(m.batch) > 14 {
		m.Emit()
	}

}

func (m *machinery) Emit() {

	m.mu.Lock()
	defer m.mu.Unlock()

	niceMap := map[string]NiceEvent{}

	for _, niceEvent := range m.batch {
		fileName := niceEvent.File
		batchedEvent, existsInMap := niceMap[fileName]

		//	only process non-temp files
		if !isTempFile(fileName) {
			if existsInMap {
				switch {
				case batchedEvent.Operation == "notify.Create" && niceEvent.Operation == "notify.Remove":
					//	a Create followed by a Remove means nothing of significance happened. Purge the record
					delete(niceMap, fileName)
				default:
					//	the default case should be to overwrite the record
					niceEvent.Topic = "rebouncer/outgoing/overwrite"
					niceMap[fileName] = niceEvent
				}
			} else {
				niceEvent.Topic = "rebouncer/outgoing/virginal"
				niceMap[niceEvent.File] = niceEvent
			}
		}

	}

	for _, ev := range niceMap {
		m.NiceChannel <- ev
	}

	m.batch = nil
	niceMap = nil

}

func (m *machinery) Subscribe() chan NiceEvent {
	return m.NiceChannel
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
