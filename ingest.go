package rebouncer

// IngestFunction eats up dirty events and produces NiceEvents.
// It _can_ decide not to push a dirty event to the the Queue
// but it _cannot_ read from the Queue. Only [ReduceFunction] can do that
// For every dirty event it consumes, it produces zero or more NiceEvents
type IngestFunction[T any] func(chan<- NiceEvent[T])

func (m *machine[T]) ingest() (chan NiceEvent[T], error) {
	incomingEvents := make(chan NiceEvent[T], m.bufferSize)
	go m.ingestor(incomingEvents)
	return incomingEvents, nil
}
