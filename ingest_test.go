package rebouncer_test

import (
	"fmt"
	"time"

	"github.com/sean9999/rebouncer"
	"github.com/rjeczalik/notify"
)


// Using NewNiceEvent in an Ingestor
func ExampleIngestFunction() {

	//	let's say we want to represent filesystem events with this struct
	type fsEvent struct {
		File          string
		Operation     string
		TransactionId uint32
	}
	
	//	this function converts from messy notify.EventInfo to our lovely fsEvent
	InotifyToNice := func(ei notify.EventInfo) rebouncer.NiceEvent[fsEvent] {
		d := fsEvent{
				File:          ei.Path(),
				Operation:     ei.Event().String(),
				TransactionId: ei.Sys().(*unix.InotifyEvent).Cookie,
		}
		r := rebouncer.NewNiceEvent[fsEvent](d, "inotify/example")
		return r
	}
	
	
	//	our IngestFunction could look like this
	ingestInotifyEvents := func(inEvents chan<- rebouncer.NiceEvent[fsEvent]) {
		//	dirty events
		var fsEvents = make(chan notify.EventInfo, 1024)
		err := notify.Watch("/var/log/...", fsEvents, WatchMask)
		if err != nil {
			panic(err)
		}
		//	clean events
		for dirtyEvent := range fsEvents {
			cleanEvent := InotifyToNice(dirtyEvent)
			//	inEvents represents the Queue
			inEvents <- cleanEvent
		}
	}	
	
}
