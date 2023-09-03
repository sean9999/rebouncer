package rebouncer

// The Emit lifecycle event happens when a value is sent to readyChannel

type EmitFunction[NICE any] func(NICE) NICE

func (m *stateMachine[NICE]) emit(fn EmitFunction[NICE]) {
	//m.Lock()
	for _, niceEvent := range m.queue {
		m.outgoingEvents <- fn(niceEvent)
	}
	//m.queue = []NICE{}
	//m.Unlock()

	m.writeQueue([]NICE{})

	if m.lifeCycleState < Draining {
		m.lifeCycleState = Emiting
	}

	//	if we are draining and the queue is empty, no need to trigger quantizer
	if m.lifeCycleState < Drained {
		if len(m.readQueue()) > 0 {
			go func() { m.readyChannel <- false }()
		} else {
			if len(m.outgoingEvents) == 0 {
				close(m.outgoingEvents)
				m.outgoingEvents = nil
				m.lifeCycleState = Drained
			}
		}
	}

}
