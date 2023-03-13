package rebouncer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
	"golang.org/x/sys/unix"
)

const WatchMask = notify.InModify | notify.InCloseWrite |
	notify.InMovedFrom | notify.InMovedTo | notify.InCreate |
	notify.InDelete | notify.InDeleteSelf | notify.InMoveSelf

type NiceEvent struct {
	Type   string
	File   string
	Event  string
	Cookie uint32
	Data   *unix.InotifyEvent
}

func (e NiceEvent) Dump() string {
	return fmt.Sprintf("%s:\t%s\t%+v (cookie=%x)", e.Event, e.File, e.Data, e.Cookie)
}

func NotifyEventInfoToNiceEvent(ei notify.EventInfo, path string, niceChannel chan NiceEvent) {
	abs, _ := filepath.Abs(path)

	data := ei.Sys().(*unix.InotifyEvent)

	n := NiceEvent{
		Type:   "fs/inotify",
		File:   strings.TrimPrefix(ei.Path(), abs+"/"),
		Event:  ei.Event().String(),
		Cookie: data.Cookie,
		Data:   ei.Sys().(*unix.InotifyEvent),
	}
	niceChannel <- n
}

func NiceEventToRebounceEvent(e NiceEvent, rbChannel chan NiceEvent) {
	e.Type = "Nice_to_rebounce"
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

	if err == nil {
		go func() {
			for fsEvent := range fsEvents {
				NotifyEventInfoToNiceEvent(fsEvent, path, niceEvents)
			}
		}()
	}

	return err
}
