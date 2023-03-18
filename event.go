package rebouncer

import (
	"fmt"
	"sync/atomic"
	"time"
)

var UniqueEventId uint64

// NiceEvent is the common format expected and produced for all events
type NiceEvent struct {
	id            uint64 // a unique (within this process) auto-incrementing ID
	TransactionId uint32 // 0 if this is an atomic operation, some number of it's part of a transaction
	Topic         string // a type of message (ex: "rebouncer/fs/inotify", or "rebouncer/fs/nice", or "rebouncer/lifecycle/shutdown")
	File          string // the file being operated on. A path relative to *watchDir
	OccurredAt    time.Time
	Operation     string // ex: Create, Delete, Modify
}

type EventMap map[string]NiceEvent

func (e NiceEvent) IsZeroed() bool {
	return (e.id == 0 && e.TransactionId == 0 && e.Topic == "")
}

func NextEventId() uint64 {
	atomic.AddUint64(&UniqueEventId, 1)
	return atomic.LoadUint64(&UniqueEventId)
}

func NewNiceEvent(topic string) NiceEvent {
	e := NiceEvent{
		id:         NextEventId(),
		Topic:      topic,
		OccurredAt: time.Now(),
	}
	return e
}

func (e NiceEvent) Dump() string {
	return fmt.Sprintf("%+v", e)
}
