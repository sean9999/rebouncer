package rebouncer

// ReduceFunction takes a slice of NiceEvents, cleans them, and returns a new slice

type ReduceFunction[T any] func([]NiceEvent[T]) []NiceEvent[T]

func (m *machine[T]) reduce(newEvent NiceEvent[T]) {
	// newEvent is added to the Queue and ReduceFunction operates on the resulting slice
	oldQueue := m.readQueue()
	newQueue := m.reducer(append(oldQueue, newEvent))

	//oldQueueLength := len(oldQueue)
	//newQueueLength := len(newQueue)
	//fmt.Printf("oldQueueLength: %d\tnewQueueLength: %d\n", oldQueueLength, newQueueLength)

	_ = m.writeQueue(newQueue)
	//fmt.Println(err)
}
