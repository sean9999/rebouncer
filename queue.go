package rebouncer

type Queue[NICE any] []NICE

func (m *stateMachine[NICE]) writeQueue(newQueue []NICE) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queue = newQueue
}

func (m *stateMachine[NICE]) readQueue() []NICE {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.queue
}
