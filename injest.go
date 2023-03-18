package rebouncer

import "github.com/rjeczalik/notify"

// listens on a channel for events, formats those events as NiceEvents, and send them along as such
type Injestor func() chan NiceEvent

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
