package rebouncer

import (
	"fmt"
	"sync/atomic"
	"time"
)

// an ever-incrementing id for NiceEvents
var UniqueEventId uint64 = uint64(time.Now().UnixNano())

// NiceEvent is the common format expected and produced for all events
type NiceEvent struct {
	Id            uint64 // a unique (within this process) auto-incrementing ID
	TransactionId uint32 // 0 if this is an atomic operation, some number of it's part of a transaction
	Topic         string // a type of message (ex: "rebouncer/fs/inotify", or "rebouncer/fs/nice", or "rebouncer/lifecycle/shutdown")
	File          string // the file being operated on. A path relative to *watchDir
	OccurredAt    time.Time
	Operation     string // ex: Create, Delete, Modify
}

// A quick check to determine of a NiceEvent is zeroed-out. Zeroed out events should be culled, but this doesn't necisarily represent an error condition.
func (e NiceEvent) IsZeroed() bool {
	return (e.Id == 0 && e.TransactionId == 0 && e.Topic == "")
}

// using atomic operations, get the next event id
func NextEventId() uint64 {
	atomic.AddUint64(&UniqueEventId, 1)
	return atomic.LoadUint64(&UniqueEventId)
}

// the canonical way to create a new event, garuanteeing you have a unique id and proper timestamp
func NewNiceEvent(topic string) NiceEvent {
	e := NiceEvent{
		Id:         NextEventId(),
		Topic:      topic,
		OccurredAt: time.Now(),
	}
	return e
}

// a simple convenience function for debugging
func (e NiceEvent) Dump() string {
	return fmt.Sprintf("%+v", e)
}
