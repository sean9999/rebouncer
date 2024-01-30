package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sean9999/GoFunctional/fslice"
	"github.com/sean9999/rebouncer"
)

func main() {

	ingestFn := func(incoming chan<- fsnotify.Event) {
		r2, err := NewRecursiveWatcher("./tmp")
		if err != nil {
			log.Fatal(err)
		}
		defer r2.Shutdown()
		for ev := range r2.Events {
			incoming <- ev
		}
	}

	reduceFn := func(oldEvents []fsnotify.Event) []fsnotify.Event {
		newEvents := make([]fsnotify.Event, 0, len(oldEvents))
		for _, thisEvent := range oldEvents {
			relatedEvents := fslice.From[fsnotify.Event](newEvents).Filter(func(e fsnotify.Event, j int, _ []fsnotify.Event) bool {
				return (e.Name == thisEvent.Name)
			})
			lastEvent := relatedEvents[len(relatedEvents)-1]
			//	for now, just include the last one
			newEvents = append(newEvents, lastEvent)
		}
		//	@todo: if we have a create followed by a write, just use create
		//	@todo: if we have a create followed by a delete, remove them both
		//	@todo: if we have anything else followed by a delete, just delete
		//	@todo: think through these and other scenarios
		return newEvents
	}

	quantizeFn := func(queue []fsnotify.Event) bool {
		//	the very moment one thing is cued, our 1-second batch is started
		//	this is a little more elegant than blindly going once a second
		//	and less likely to chop a batch in half

		if len(queue) == 0 {
			return false
		} else {
			time.Sleep(time.Millisecond * 1000)
			return true
		}
	}

	rebby := rebouncer.NewRebouncer[fsnotify.Event](
		ingestFn,
		reduceFn,
		quantizeFn,
		1024,
	)

	for ev := range rebby.Subscribe() {
		fmt.Println(ev)
	}

}
