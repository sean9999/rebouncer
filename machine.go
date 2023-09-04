package rebouncer

import (
	"fmt"
	"sync"
)

type lifeCycleState int

const (
	StartingUp lifeCycleState = iota
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
	lifeCycle      chan lifeCycleState
	incomingEvents chan NICE
	outgoingEvents chan NICE
	queue          Queue[NICE]
	mu             sync.RWMutex
	lifeState      lifeCycleState
}

func (m *stateMachine[NICE]) SetLifeCycleState(s lifeCycleState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lifeState = s
	fmt.Println(s)
}

func (m *stateMachine[NICE]) GetLifeCycleState() lifeCycleState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lifeState
}

func (m *stateMachine[NICE]) Interrupt() {
	m.lifeCycle <- Draining
}

func (m *stateMachine[NICE]) Subscribe() <-chan NICE {
	return m.outgoingEvents
}
