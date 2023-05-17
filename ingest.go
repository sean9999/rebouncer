package rebouncer

// IngestFunction eats up dirty events and produces NiceEvents.
// It can decide to simply convert all dirty events to their clean equivilants,
// or drop some on the floor.
// 
// An IngestFunction can only push new NiceEvents to the queue. It doesn't know what's already there. 
// Ingest is the first lifecycle event. It will be followed by [ReduceFunction]
type IngestFunction[T any] func(chan<- NiceEvent[T])

func (m *machine[T]) ingest() (chan NiceEvent[T], error) {
	incomingEvents := make(chan NiceEvent[T], m.bufferSize)
	go m.ingestor(incomingEvents)
	return incomingEvents, nil
}
