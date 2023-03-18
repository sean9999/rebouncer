package rebouncer

import "github.com/rjeczalik/notify"

// listens on a channel for events, formats those events as NiceEvents, and send them along to a chanel it returns
type Injestor func() chan NiceEvent

// DefaultInotifyInjestor satisfies the [Injestor] type, listening for inotify events on a directory, formatting those events, and sending them along to a channel
func DefaultInotifyInjestor(dir string, bufferSize int) Injestor {
	var niceChan = make(chan NiceEvent, bufferSize)
	var fsEvents = make(chan notify.EventInfo, bufferSize)
	jester := func() chan NiceEvent {
		err := notify.Watch(dir+"/...", fsEvents, WatchMask)
		if err != nil {
			panic(err)
		}
		go func() {
			for fsEvent := range fsEvents {
				niceChan <- NotifyToNiceEvent(fsEvent, dir)
			}
		}()
		return niceChan
	}
	return jester
}
