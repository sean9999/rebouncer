package rebouncer_test

import (
	"fmt"
	"path/filepath"

	"github.com/sean9999/rebouncer"
)

// A Reducer operates on the Queue of NiceEvents waiting to be flushed
func ExampleReduceFunction() {

	//	let's say this is our underyling filesystem event
	type fsEvent struct {
		File          string
		Operation     string
		TransactionId uint32
	}

	//	we wrap it in a NiceEvent, because rebouncer always works with NiceEvents
	type MyNiceEvent rebouncer.NiceEvent[fsEvent]
	
	//	omitCss is our ReduceFunction 
	omitCss := func(inEvents []MyNiceEvent) []MyNiceEvent {
		outEvents := []MyNiceEvent{}
		//	omit any event on a CSS file, but include all others
		//	a Reducer can modify events too
		for _, e := range inEvents {
			if filepath.Ext(e.Data.File) != ".css" {
				e.Topic = "inotify/normalized"
				outEvents = append(outEvents, e)
			}
		}
		return outEvents
	}

}
