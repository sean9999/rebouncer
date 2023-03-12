package pkg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

const BUFFER_LENGTH = 1024

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

// WatchRecursively emits event info to the "niceEvents" channel
func WatchRecursively(path string, niceEvents chan NiceEvent) error {

	var fsEvents = make(chan notify.EventInfo, BUFFER_LENGTH)
	err := notify.Watch(path+"/...", fsEvents, notify.All)

	rebouncedEvents := make(chan NiceEvent, BUFFER_LENGTH)

	if err == nil {

		go func() {
			for {
				select {
				case niceEvent := <-niceEvents:
					fmt.Println("niceEvent", niceEvent)
					NiceEventToRebounceEvent(niceEvent, rebouncedEvents)
				case rb := <-rebouncedEvents:
					fmt.Println("rb", rb)
				case fsEvent := <-fsEvents:
					fmt.Println("fsEvent", fsEvent)
					NotifyEventInfoToNiceEvent(fsEvent, path, niceEvents)
				}
			}
		}()
	}

	return err
}
