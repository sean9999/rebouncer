package rebouncer

// takes an array of NiceEvents, filters out undesired ones, and returns a cleaner array of NiceEvents
type Reducer func([]NiceEvent) []NiceEvent

func DefaultInotifyReducer(inEvents []NiceEvent) []NiceEvent {

	batchMap := EventMap{}
	normalizedEvents := []NiceEvent{}

	//	fill batchMap
	for _, newEvent := range inEvents {
		fileName := newEvent.File
		oldEvent, existsInMap := batchMap[fileName]
		//	only process non-temp files and valid events
		if !isTempFile(fileName) && !newEvent.IsZeroed() {
			if existsInMap {
				switch {
				case oldEvent.Operation == "notify.Create" && newEvent.Operation == "notify.Remove":
					//	a Create followed by a Remove means nothing of significance happened. Purge the record
					delete(batchMap, fileName)
				default:
					//	the default case should be to overwrite the record
					newEvent.Topic = "rebouncer/outgoing/1"
					batchMap[fileName] = newEvent
				}
			} else {
				newEvent.Topic = "rebouncer/outgoing/0"
				batchMap[newEvent.File] = newEvent
			}
		}
	}

	//	unwind batchMap
	for _, e := range batchMap {
		normalizedEvents = append(normalizedEvents, e)
	}

	return normalizedEvents

}
