package rebouncer

// The Emit lifecycle event happens when a value is sent to readyChannel

func (m *rebounceMachine[NICE]) drainQueue() {
	for _, niceEvent := range m.queue {
		m.outgoingEvents <- niceEvent
	}
	m.writeQueue([]NICE{})
}

func (m *rebounceMachine[NICE]) emit() {

	m.drainQueue()

	if m.GetLifeCycleState() == Draining {
		m.lifeCycle <- Drained
	} else {
		m.lifeCycle <- Quantizing
	}

}
