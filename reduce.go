package rebouncer

// takes an array of NiceEvents, filters out undesired ones, and returns a cleaner array of NiceEvents
type Reducer func([]NiceEvent) []NiceEvent

// DefaultInotifyReduce is a convenience function satisfying type [Reducer], performing the following cleanup:
//   - removes all events relating to temp files
//   - a create followed by a write on the same file becomes just one event (a create)
//   - a delete followed by a create on the same underlying file descriptor becaomes just one rename
//   - in the case of a create followed by a delete on the same file, both are removed.
func DefaultInotifyReduce(inEvents []NiceEvent) []NiceEvent {

	// folding our slice into a map and then back into a slice is a convenient way to normalize
	// because we are using FileName as key
	batchMap := map[string]NiceEvent{}
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
					newEvent.Topic = "rebouncer/inotify/outgoing/clobber"
					batchMap[fileName] = newEvent
				}
			} else {
				newEvent.Topic = "rebouncer/inotify/outgoing/virginal"
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
