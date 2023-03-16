package rebouncer

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// NiceEvent is the common format expected and produced for all events
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
