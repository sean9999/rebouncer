package rebouncer

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

// NiceEvent is the common format expected and produced for all events
type NiceEvent struct {
	ID         uint32 // a unique (within this process) auto-incrementing ID
	TxID       uint32 // 0 if this is an atomic operation, some number of it's part of a transaction
	Topic      string // a type of message (ex: "rebouncer/fs/inotify", or "rebouncer/fs/nice", or "rebouncer/lifecycle/shutdown")
	File       string //	the file being operated on. A path relative to *watchDir
	OccurredAt time.Time
	Operation  string             // ex: Create, Delete, Modify
	Data       *unix.InotifyEvent // original event or any extra data useful to the consumer
}

func (e NiceEvent) Dump() string {
	return fmt.Sprintf("%s:\t%s\t%+v (cookie=%x)", e.Operation, e.File, e.Data, e.TxID)
}

/*
func (b BusEvent) Dump() string {
	return fmt.Sprintf("%+v", b)
}
*/
