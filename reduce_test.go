package rebouncer_test

import (
	"path/filepath"

	"github.com/sean9999/rebouncer"
)

func ExampleReduceFunction() {

	//	a filesystem event
	type fsEvent struct {
		File          string
		Operation     string
		TransactionId uint64
	}

	//	we wrap it in a NiceEvent, because rebouncer always expects this
	type MyNiceEvent rebouncer.NiceEvent[fsEvent]

	//	ReduceFunction[fsEvent]
	_ = func(inEvents []MyNiceEvent) []MyNiceEvent {
		outEvents := []MyNiceEvent{}
		//	omit any event on a CSS file, but include all others
		//	set the topic on all to "inotify/normalized"
		for _, e := range inEvents {
			if filepath.Ext(e.Data.File) != ".css" {
				e.Topic = "inotify/normalized"
				outEvents = append(outEvents, e)
			}
		}
		return outEvents
	}

}
