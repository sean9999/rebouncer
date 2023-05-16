package rebouncer

type Behaviour[T any] interface {
	Subscribe() <-chan NiceEvent[T]     // returns a channel and pushes events to it
	emit()                              // flushes the Queue
	readQueue() []NiceEvent[T]          // gets the Queue, with safety and locking
	writeQueue([]NiceEvent[T]) error    // sets the Queue, handling safety and locking
	ingest() (chan NiceEvent[T], error) // returns and operates on a channel
	quantize()                          // decides whether the flush the Queue
	reduce(NiceEvent[T])                // removes unwanted NiceEvents from the Queue
}
