package rebouncer

// IngestFunction eats up dirty events and produces NiceEvents.
// It can decide to simply convert all dirty events to their clean equivalents,
// or drop some on the floor.
//
// An IngestFunction can only push new NiceEvents to the queue. It doesn't know what's already there.
// Ingest is the first lifecycle event. It will be followed by [ReduceFunction]
type IngestFunction[NICE any] func(chan<- NICE)

func (m *stateMachine[NICE]) ingest(fn IngestFunction[NICE]) {

	go func() {
		fn(m.incomingEvents)    // run ingest function to completion
		close(m.incomingEvents) // when it's done, that means there are no more incoming events
		m.lifeCycle <- Draining // therefore, we're draining
	}()

}
