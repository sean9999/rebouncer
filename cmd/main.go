package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/sean9999/rebouncer"
)

var watchDir *string
var flushPeriod *int

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	flushPeriod = flag.Int("period", 1000, "how often (in milliseconds) to flush events")
	flag.Parse()
}

/*

{	id:1679179339846315599
	TransactionId:0
	Topic:rebouncer/inotify/outgoing/clobber
	File:f2.txt
	OccurredAt:2023-03-18 18:42:33.557967981 -0400 EDT m=+13.711725019
	Operation:notify.InCloseWrite
}
{
	id:1679179339846315605
	TransactionId:0
	Topic:rebouncer/inotify/outgoing/clobber
	File:f3.txt
	OccurredAt:2023-03-18 18:56:25.787518001 -0400 EDT m=+845.941275029
	Operation:notify.InCloseWrite
}

event: rebouncer/fs/output
data: {"file": "index.html", "operation": "modify"}

event: rebouncer/fs/output
data: {"file": "css/debug.css", "operation": "delete"}

event: rebouncer/fs/output
data: {"file": "css/mobile", "operation": "create"}

event: rebouncer/fs/output
data: {"file": "css/mobile", "operation": "modify"}

*/

type fsEvent struct {
	File      string `json:"file"`
	Operation string `json:"operation"`
}

type SSEEvent struct {
	Id    uint64  `json:"id"`
	Event string  `json:"event"`
	Data  fsEvent `json:"data"`
}

func (s SSEEvent) Serialize() string {
	eventDataAsJson, err := json.Marshal(s.Data)
	if err != nil {
		panic(err)
	}
	output := fmt.Sprintf("event: %s\nid: %d\ndata: %s\n\n", s.Event, s.Id, eventDataAsJson)
	return output
}

func NiceEventToSSE(ne rebouncer.NiceEvent) SSEEvent {
	inotifyEvent := fsEvent{
		File:      ne.File,
		Operation: ne.Operation,
	}
	sseEvent := SSEEvent{
		Id:    ne.Id,
		Event: ne.Topic,
		Data:  inotifyEvent,
	}
	return sseEvent
}

func main() {

	//	instantiate
	niceEvents := rebouncer.NewInotify(*watchDir, *flushPeriod)

	//	here is the channel we listen on
	outgoingEvents := niceEvents.Subscribe()

	//	output in SSE format
	for e := range outgoingEvents {
		fmt.Printf("%s", NiceEventToSSE(e).Serialize())
	}

}
