package rebouncer

type Queue[NICE any] []NICE

func (m *stateMachine[NICE]) Lock() {
	m.mu.Lock()
}

func (m *stateMachine[NICE]) Unlock() {
	m.mu.Lock()
}

func (m *stateMachine[NICE]) writeQueue(newQueue []NICE) {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	m.queue = newQueue
}

func (m *stateMachine[NICE]) readQueue() []NICE {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	return m.queue
}
