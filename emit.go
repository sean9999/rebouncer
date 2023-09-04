package rebouncer

// The Emit lifecycle event happens when a value is sent to readyChannel

type EmitFunction[NICE any] func(NICE) NICE

func (m *stateMachine[NICE]) drainQueue(fn EmitFunction[NICE]) {
	for _, niceEvent := range m.queue {
		m.outgoingEvents <- fn(niceEvent)
	}
	m.writeQueue([]NICE{})
}

func (m *stateMachine[NICE]) emit(fn EmitFunction[NICE]) {

	m.drainQueue(fn)

	if m.GetLifeCycleState() == Draining {
		m.lifeCycle <- Drained
	} else {
		m.lifeCycle <- Quantizing
	}

}
