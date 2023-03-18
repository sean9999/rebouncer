package rebouncer

import (
	"fmt"
	"time"
)

// Quantizer runs in a go routine and emits to machinery.ticker.C when it decides we're ready to batch up our events
type Quantizer func(chan bool, *[]NiceEvent)

// simply waits ms milliseconds and then sends true, causing Emit() to run, sending NiceEvents back to the consumer
// this would be the most common and straightforward pattern for filesystem watchers
func DefaultInotifyQuantizer(ms int) Quantizer {
	ticker := time.NewTicker(time.Minute)
	qFunc := func(readyChan chan bool, em *[]NiceEvent) {
		fmt.Println("quantizer func", ms)
		ticker.Reset(3 * time.Second)
		for range ticker.C {
			lengthOfArray := len(*em)
			ready := (len(*em) > 0)
			fmt.Println("ticker", ready, lengthOfArray)
			readyChan <- ready
		}
	}
	return qFunc
}
