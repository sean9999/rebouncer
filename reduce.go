package rebouncer

// Reducer modifies the [Queue]. It takes a slice of events, cleans them, and returns a new slice.
// [Reducing] is the 2nd lifecycle event, after [Ingesting] and "before" [Quantizing].
//
// [Quantizer] actually runs in it's own loop seperate from ingest=>reduce, but it's helpful to think of it as coming after Reduce.
type Reducer[NICE any] func([]NICE) []NICE

func (m *rebounceMachine[NICE]) reduce(fn Reducer[NICE], newEvent NICE) {
	//	apply ReduceFunction to the queue with the new event appended
	//	write the result back to the queue
	newQueue := fn(append(m.readQueue(), newEvent))
	m.writeQueue(newQueue)
}
