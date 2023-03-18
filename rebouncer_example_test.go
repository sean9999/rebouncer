package rebouncer_test

import (
	"fmt"

	"github.com/sean9999/rebouncer"
)

// This is a basic example showing how to work with notify.Watch function.
func ExampleNew() {

	interval := 3000      // how long to wait in between flushes, in milliseconds
	watchDir := "./build" // what directory to recursively watch for fileSystem events

	//	rebecca is our singleton instance
	rebecca := rebouncer.New(rebouncer.Config{
		BufferSize: rebouncer.DefaultBufferSize,
		Quantizer:  rebouncer.DefaultInotifyQuantizer(interval),
		Reducer:    rebouncer.DefaultInotifyReducer,
		Injestor:   rebouncer.DefaultInotifyInjestor(watchDir, rebouncer.DefaultBufferSize),
	})

	//	here is the channel we can listen on
	outgoingEvents := rebecca.Subscribe()

	//	for example
	for e := range outgoingEvents {
		fmt.Println(e.Dump())
	}

}
