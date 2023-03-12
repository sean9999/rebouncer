package rebouncer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

// formats inotify events to our liking
func NotifyEventInfoToNiceEvent(ei notify.EventInfo, path string) NiceEvent {

	fmt.Println("notifyEventInfoToNiceEvent")

	abs, _ := filepath.Abs(path)
	return NiceEvent{
		File:  strings.TrimPrefix(ei.Path(), abs+"/"),
		Event: ei.Event().String(),
	}
}

// Watch emits inotify events to the "niceEvents" channel, basically just formatting them
func Watch(rootPath string, niceEvents chan NiceEvent) error {

	fmt.Println("Watch")

	var fsEvents = make(chan notify.EventInfo)
	err := notify.Watch(rootPath+"/...", fsEvents, notify.All)

	//	massage the event to the format we want
	func() {
		for eventInfo := range fsEvents {
			niceEvent := NotifyEventInfoToNiceEvent(eventInfo, rootPath)
			niceEvents <- niceEvent
		}
	}()

	return err
}
