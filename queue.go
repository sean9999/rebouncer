package rebouncer

type Queue[NICE any] []NICE

func (m *rebounceMachine[NICE]) writeQueue(newQueue []NICE) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queue = newQueue
}

func (m *rebounceMachine[NICE]) readQueue() []NICE {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.queue
}
