package rebouncer

// IngestFunction eats up dirty events and produces NiceEvents.
// It can decide to simply convert all dirty events to their clean equivalents,
// or drop some on the floor.
//
// An IngestFunction can only push new NiceEvents to the queue. It doesn't know what's already there.
// Ingest is the first lifecycle event. It will be followed by [ReduceFunction]
type IngestFunction[NICE any] func(chan<- NICE, chan bool)

func (m *stateMachine[NICE]) ingest(fn IngestFunction[NICE]) {
	go fn(m.incomingEvents, m.doneChannel)
}

func (m *stateMachine[NICE]) Drain() {
	m.lifeCycleState = Draining
	close(m.incomingEvents)
}
