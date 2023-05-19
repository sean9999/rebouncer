package rebouncer

type Queue[NICE any] []NICE

func (m *stateMachine[NAUGHTY, NICE, BEAUTIFUL]) writeQueue(newQueue []NICE) {
	m.mu.Lock()
	m.queue = newQueue
	m.mu.Unlock()
}

func (m *stateMachine[NAUGHTY, NICE, BEAUTIFUL]) readQueue() []NICE {
	return m.queue
}
