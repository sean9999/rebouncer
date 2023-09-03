package rebouncer

// ReduceFunction takes a slice of NiceEvents, cleans them, and returns a new slice
// It was designed with the use-case of removing extraneous events, but it could just as easily add new ones
// or modify existing ones. A [ReduceFunction] is the 2nd lifecycle event, after [IngestFunction] and before [QuantizeFunction]
type ReduceFunction[NICE any] func([]NICE) []NICE

func (m *stateMachine[NICE]) reduce(fn ReduceFunction[NICE], newEvent NICE) {
	//	apply ReduceFunction to the queue with the new NICE event appended
	//	write the result back to the queue
	if m.lifeCycleState < Draining {
		m.lifeCycleState = Reducing
	} else {
		if len(m.incomingEvents) == 0 {
			close(m.incomingEvents)
			m.incomingEvents = nil
		}
	}
	newQueue := fn(append(m.readQueue(), newEvent))
	m.writeQueue(newQueue)
}
