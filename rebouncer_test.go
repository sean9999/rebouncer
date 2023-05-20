package rebouncer_test

import (
	"fmt"
	"testing"

	"github.com/sean9999/rebouncer"
)

type CanonicalDirtyEvent struct {
	Id           string
	UsefulField  string
	UselessField string
}

type CanonicalNiceEvent struct {
	CanonicalDirtyEvent
	timeStamp int64
}

type CanonicalBeautifulEvent struct {
	id          string
	usefulField string
}

var CanonicalIngestFunction rebouncer.IngestFunction[CanonicalDirtyEvent, CanonicalNiceEvent] = func(chan<- CanonicalNiceEvent) {

	dirtyEvents := []CanonicalDirtyEvent{}

	for i := 0; i < 100; i++ {
		thisEvent := CanonicalDirtyEvent{
			Id:           fmt.Sprintf("%d", i),
			UsefulField:  fmt.Sprintf("I am record %d", i),
			UselessField: fmt.Sprintf("I am also record %d", i),
		}
		dirtyEvents = append(dirtyEvents, thisEvent)
	}

}

func TestNewRebouncer(t *testing.T) {

	t.Run("create a rebouncer with three structs and no user-defined functions", func(t *testing.T) {

		rebecca := rebouncer.NewRebouncer[dirtyEvent, niceEvent, beautifulEvent]()
	})

}
