package rebouncer

import (
	"sync"
)

type rebouncerLifecycleState int

const (
	StartingUp rebouncerLifecycleState = iota
	Running
	Ingesting
	Reducing
	Quantizing
	Emiting
	Draining
	Drained
	ShuttingDown
)

// *stateMachine implements [Behaviour] and contains state
type stateMachine[NICE any] struct {
	//config         Config
	//user           UserDefinedFunctionSet[NAUGHTY, NICE, BEAUTIFUL]
	readyChannel   chan bool
	doneChannel    chan bool // to indicate we're done ingesting
	incomingEvents chan NICE
	outgoingEvents chan NICE
	queue          Queue[NICE]
	lifeCycleState rebouncerLifecycleState
	mu             sync.Mutex
}

func (m *stateMachine[NICE]) Done() {
	m.doneChannel <- true
}

func (m *stateMachine[NICE]) Subscribe() <-chan NICE {
	return m.outgoingEvents
}
