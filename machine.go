package rebouncer

import (
	"sync"
)

// *stateMachine implements [Behaviour] and contains state
type stateMachine[NAUGHTY any, NICE any, BEAUTIFUL any] struct {
	//config         Config
	//user           UserDefinedFunctionSet[NAUGHTY, NICE, BEAUTIFUL]
	readyChannel   chan bool
	incomingEvents chan NICE
	outgoingEvents chan BEAUTIFUL
	queue          Queue[NICE]
	mu             sync.Mutex
}

func (m *stateMachine[NAUGHTY, NICE, BEAUTIFUL]) Subscribe() <-chan BEAUTIFUL {
	return m.outgoingEvents
}
