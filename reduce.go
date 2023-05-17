package rebouncer

// ReduceFunction takes a slice of NiceEvents, cleans them, and returns a new slice
// It was designed with the use-case of removing extraneous events, but it could just as easily add new ones
// or modify existing ones. A [ReduceFunction] is the 2nd lifecycle event, after [IngestFunction] and before [QuantizeFunction]
type ReduceFunction[T any] func([]NiceEvent[T]) []NiceEvent[T]

func (m *machine[T]) reduce(newEvent NiceEvent[T]) {
	// newEvent is added to the Queue. ReduceFunction operates on the resulting slice
	oldQueue := m.readQueue()
	newQueue := m.reducer(append(oldQueue, newEvent))

	_ = m.writeQueue(newQueue)
}
