package rebouncer

import (
	"time"
)

// Quantizer runs in a go routine and sends true to readyChannel when it decides we're ready to Emit()
// it has access to the entire batch in the Queue to help make this decision.
type Quantizer func(chan bool, *[]NiceEvent)

// simply waits ms milliseconds and then sends true, causing Emit() to run, sending NiceEvents back to the consumer
// this would be the most common and straightforward pattern for filesystem watchers
func DefaultInotifyQuantizer(ms int) Quantizer {
	ticker := time.NewTicker(time.Minute)
	qFunc := func(readyChan chan bool, em *[]NiceEvent) {
		ticker.Reset(time.Duration(ms) * time.Millisecond)
		for range ticker.C {
			ready := (len(*em) > 0)
			readyChan <- ready
		}
	}
	return qFunc
}
