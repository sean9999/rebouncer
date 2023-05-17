package rebouncer_test

import (
	"fmt"
	"time"

	"github.com/sean9999/rebouncer"
	"github.com/rjeczalik/notify"
)


// Using NewNiceEvent in an Ingestor
func ExampleQuantizeFunction() {
	
	interval, _ := time.ParseDuration("1s")
	
	type MyNiceEvent rebouncer.NiceEvent[fsEvent]
	
	periodicallyFlushQueue := func(readyChan chan<- bool, queue []MyNiceEvent) {
		if len(queue) > 0 {
			readyChan <- true
		} else {
			time.Sleep(interval)
			readyChan <- false
		}
	}

}
