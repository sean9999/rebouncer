package rebouncer

// ReduceFunction takes a slice of NiceEvents, cleans them, and returns a new slice
// It was designed with the use-case of removing extraneous events, but it could just as easily add new ones
// or modify existing ones. A [ReduceFunction] is the 2nd lifecycle event, after [IngestFunction] and before [QuantizeFunction]
type ReduceFunction[NICE any] func([]NICE) []NICE

func (m *stateMachine[NAUGHTY, NICE, BEAUTIFUL]) reduce(fn ReduceFunction[NICE], newEvent NICE) {
	//	apply ReduceFunction to the queue with the new NICE event appended
	//	write the result back to the queue
	m.writeQueue(fn(append(m.readQueue(), newEvent)))
}

/*
func (m *machine[T]) reduce(newEvent NiceEvent[T]) {
	// newEvent is added to the Queue. ReduceFunction operates on the resulting slice
	oldQueue := m.readQueue()
	newQueue := m.reducer(append(oldQueue, newEvent))

	_ = m.writeQueue(newQueue)
}
*/
