package rebouncer

// IngestFunction eats up dirty events and produces NiceEvents.
// It can decide to simply convert all dirty events to their clean equivalents,
// or drop some on the floor.
//
// An IngestFunction can only push new NiceEvents to the queue. It doesn't know what's already there.
// Ingest is the first lifecycle event. It will be followed by [ReduceFunction]
type IngestFunction[NAUGHTY any, NICE any] func(chan<- NICE)

func (m *stateMachine[NAUGHTY, NICE, BEAUTIFUL]) ingest(fn IngestFunction[NAUGHTY, NICE]) {
	go fn(m.incomingEvents)
}
