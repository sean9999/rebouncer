package rebouncer_test

import (
	"fmt"
	"path/filepath"

	"github.com/sean9999/rebouncer"
)

// A Reducer operates on the Queue of NiceEvents waiting to be flushed
func ExampleReducer() {

	omitCss := func(inEvents []rebouncer.NiceEvent) []rebouncer.NiceEvent {
		outEvents := []rebouncer.NiceEvent{}
		//	omit any event on a CSS file, but inlcude all others
		//	a Reducer can modify events too
		for _, e := range inEvents {
			if filepath.Ext(e.File) != ".css" {
				e.Topic = "rebouncer/fs/normalized"
				outEvents = append(outEvents, e)
			}
		}
		return outEvents
	}

	// Reducers (as well as ingestors and Quantizers) are injected at instantiation time
	conf := rebouncer.Config{
		Reducer: omitCss,
	}
	reb := rebouncer.New(conf)

	// the resulting events will not contain CSS files
	for e := range reb.Subscribe() {
		fmt.Println(e.Dump())
	}

}
