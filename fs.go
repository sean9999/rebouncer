package rebouncer

import (
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
	"golang.org/x/sys/unix"
)

const WatchMask = notify.InModify | notify.InCloseWrite |
	notify.InMovedFrom | notify.InMovedTo | notify.InCreate |
	notify.InDelete | notify.InDeleteSelf | notify.InMoveSelf

func NotifyEventInfoToNiceEvent(ei notify.EventInfo, path string, niceChannel chan NiceEvent) {
	abs, _ := filepath.Abs(path)

	data := ei.Sys().(*unix.InotifyEvent)

	n := NiceEvent{
		Topic:     "rebouncer/fs/inotify",
		File:      strings.TrimPrefix(ei.Path(), abs+"/"),
		Operation: ei.Event().String(),
		TxID:      data.Cookie,
		Data:      data,
	}
	niceChannel <- n
}

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

	var fsEvents = make(chan notify.EventInfo, DefaultBufferSize)
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
