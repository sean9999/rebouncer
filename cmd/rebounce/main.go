package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/sean9999/rebouncer"
	"golang.org/x/sys/unix"
)

// the string ends with a "~" character
func endsInTilde(s string) bool {
	pattern := `~$`
	r, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}
	return r
}

// the string is just numbers
func containsOnlyNumbers(s string) bool {
	pattern := `^\d+$`
	r, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}
	return r
}

func isTempFile(path string) bool {
	return endsInTilde(path) && containsOnlyNumbers(path)
}

func InotifyToNice(ei notify.EventInfo, basePath string) rebouncer.NiceEvent[fsEvent] {
	abs, _ := filepath.Abs(basePath)
	normalFile := strings.TrimPrefix(ei.Path(), abs)
	d := fsEvent{
		File:          normalFile,
		Operation:     ei.Event().String(),
		TransactionId: ei.Sys().(*unix.InotifyEvent).Cookie,
	}
	r := rebouncer.NewNiceEvent[fsEvent](d, "rebouncer/inotify/example")
	return r
}

type fsEvent struct {
	File          string
	Operation     string
	TransactionId uint32
}

// Instantiating an inotify-backed file-watcher using the verbose method.
func main() {

	const WatchMask = notify.InModify |
		notify.InCloseWrite |
		notify.InMovedFrom |
		notify.InMovedTo |
		notify.InCreate |
		notify.InDelete |
		notify.InDeleteSelf |
		notify.InMoveSelf

	interval, err := time.ParseDuration("1001ms") // how long to wait in between flushes, in milliseconds
	if err != nil {
		panic(err)
	}
	watchDir := "./build" // what directory to recursively watch for fileSystem events

	ingestInotifyEvents := func(inEvents chan<- rebouncer.NiceEvent[fsEvent]) {
		//	dirty events
		var fsEvents = make(chan notify.EventInfo, 1024)
		err := notify.Watch(watchDir+"/...", fsEvents, WatchMask)
		if err != nil {
			panic(err)
		}
		// clean events
		for dirtyEvent := range fsEvents {
			cleanEvent := InotifyToNice(dirtyEvent, watchDir)
			if !isTempFile(cleanEvent.Data.File) {
				inEvents <- cleanEvent
			}
		}
	}

	removeDuplicateInotifyEvents := func(inEvents []rebouncer.NiceEvent[fsEvent]) []rebouncer.NiceEvent[fsEvent] {
		// folding our slice into a map and then back into a slice is a convenient way to normalize
		// because we are using FileName as key
		batchMap := map[string]rebouncer.NiceEvent[fsEvent]{}
		normalizedEvents := []rebouncer.NiceEvent[fsEvent]{}
		//	fill batchMap
		for _, newEvent := range inEvents {
			fileName := newEvent.Data.File
			oldEvent, existsInMap := batchMap[fileName]
			//	only process non-temp files and valid events
			//if !isTempFile(fileName) && !newEvent.IsZeroed() {
			if existsInMap {
				switch {
				case oldEvent.Data.Operation == "notify.Create" && newEvent.Data.Operation == "notify.Remove":
					//	a Create followed by a Remove means nothing of significance happened. Purge the record
					delete(batchMap, fileName)
				default:
					//	the default case should be to overwrite the record
					newEvent.Topic = "rebouncer/inotify/outgoing/clobber"
					batchMap[fileName] = newEvent
				}
			} else {
				newEvent.Topic = "rebouncer/inotify/outgoing/virginal"
				batchMap[newEvent.Data.File] = newEvent
			}
			//}
		}
		//	unwind batchMap
		for _, e := range batchMap {
			normalizedEvents = append(normalizedEvents, e)
		}

		return normalizedEvents
	}

	periodicallyFlushQueue := func(readyChan chan<- bool, queue []rebouncer.NiceEvent[fsEvent]) {
		period := time.NewTicker(interval)
		for range period.C {
			queueLength := len(queue)
			fmt.Println("queueLength", queueLength)
			if queueLength > 0 {
				readyChan <- true
			} else {
				readyChan <- false
			}
		}
	}

	conf := rebouncer.Config[fsEvent]{
		BufferSize: 1024,
		Ingestor:   ingestInotifyEvents,
		Reducer:    removeDuplicateInotifyEvents,
		Quantizer:  periodicallyFlushQueue,
	}

	//	rebecca is our singleton instance
	rebecca, err := rebouncer.New(conf)
	if err != nil {
		panic(err)
	}

	//fmt.Println(rebecca)

	//	here is the channel we can listen on
	outgoingEvents := rebecca.Subscribe()

	//	for example
	for e := range outgoingEvents {
		fmt.Println("REBOUNCER", e.Dump())
	}

}
