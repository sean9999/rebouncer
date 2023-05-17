package rebouncer

import (
	"fmt"
	"sync/atomic"
	"time"
)

// UniqueEventId is an ever-incrementing id for NiceEvents
var UniqueEventId uint64 = uint64(time.Now().UnixNano())

// NiceEvent is the common format expected and produced for all events
type NiceEvent[T any] struct {
	Data  T // original data
	Id    uint64
	Topic string
}

// NextEventId uses atomic operations to increment a [UniqueEventId]
func nextEventId() uint64 {
	atomic.AddUint64(&UniqueEventId, 1)
	return atomic.LoadUint64(&UniqueEventId)
}

// IsZeroed is a quick check to determine of a NiceEvent is zeroed-out.
// Zeroed out events should be culled, but this doesn't necessarily represent an error condition.
// Your business logic will depend on your needs
func (e NiceEvent[T]) IsZeroed() bool {
	return e.Id == 0 && e.Topic == ""
}

// NewNiceEvent is the canonical way to create a new event, guaranteeing you have a unique id and proper timestamp
func NewNiceEvent[T any](originalEvent T, topic string) NiceEvent[T] {
	e := NiceEvent[T]{
		Id:    nextEventId(),
		Topic: topic,
		Data:  originalEvent,
	}
	return e
}

// Dump is a simple convenience function for debugging
func (e NiceEvent[T]) Dump() string {
	return fmt.Sprintf("NiceEvent: %+v", e)
}
