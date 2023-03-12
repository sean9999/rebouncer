package rebouncer

// takes an array of NiceEvents, filters out the useless ones, and returns another array of NiceEvents
func Reduce(batchedEvents chan []NiceEvent, rebouncedEvents chan NiceEvent) {

	normalizedEvents := []NiceEvent{}
	b := <-batchedEvents

	for _, e := range b {
		//	do some logic
		if e.File != "." {
			normalizedEvents = append(normalizedEvents, e)
		}
	}

	for _, e := range normalizedEvents {
		rebouncedEvents <- e
	}

}
