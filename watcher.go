package rebouncer

import (
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

type NiceEvent struct {
	Tag   string
	File  string
	Event string
}

func NotifyEventInfoToNiceEvent(ei notify.EventInfo, path string, niceChannel chan NiceEvent) {
	abs, _ := filepath.Abs(path)
	n := NiceEvent{
		Tag:   "EventInfo_to_Nice",
		File:  strings.TrimPrefix(ei.Path(), abs+"/"),
		Event: ei.Event().String(),
	}
	niceChannel <- n
}

func NiceEventToRebounceEvent(e NiceEvent, rbChannel chan NiceEvent) {
	e.Tag = "Nice_to_rebounce"
	rbChannel <- e
}

/*
func thisIsTheLastEventInThisArrayReferencingThisFilename(j int, eventInQuestion NiceEvent, arr []NiceEvent) bool {

	var x int

	for i, e := range arr {
		if e.File == eventInQuestion.File {
			x = i
		}
	}

	return (j == x)
}
*/

func NormalizeEvents(inEvents []NiceEvent) []NiceEvent {
	var r []NiceEvent = inEvents

	/*
		for i, thisEvent := range inEvents {
			if thisIsTheLastEventInThisArrayReferencingThisFilename(i, thisEvent, inEvents) {
				r = append(r, thisEvent)
			}
		}
	*/

	return r
}

// WatchDirectory emits events to the "niceEvents" channel
func WatchDirectory(path string, niceEvents chan NiceEvent) error {

	var fsEvents = make(chan notify.EventInfo, BufferLength)
	err := notify.Watch(path+"/...", fsEvents, notify.All)

	var batchedEvents []NiceEvent

	if err == nil {
		go func() {
			for {
				select {
				case niceEvent := <-niceEvents:
					//fmt.Println("niceEvent", niceEvent)
					//NiceEventToRebounceEvent(niceEvent, rebouncedEvents)
					batchedEvents = append(batchedEvents, niceEvent)
				case fsEvent := <-fsEvents:
					//fmt.Println("fsEvent", fsEvent)
					NotifyEventInfoToNiceEvent(fsEvent, path, niceEvents)
					/*
						case <-time.After(3 * time.Second):
							if len(batchedEvents) > 0 {
								normalizedEvents := batchedEvents
								fmt.Println("batch")
								for _, e := range normalizedEvents {
									niceEvents <- e
								}
								batchedEvents = nil
							}
					*/
				}

			}
		}()
	}

	return err
}
