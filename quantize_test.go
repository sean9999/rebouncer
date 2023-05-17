package rebouncer_test

import (
	"time"

	"github.com/sean9999/rebouncer"
)

func ExampleQuantizeFunction() {

	type fsEvent struct {
		File          string
		Operation     string
		TransactionId uint64
	}

	type MyNiceEvent rebouncer.NiceEvent[fsEvent]

	//	since QuantizeFunction is a closure, it can access its outer scope
	interval, _ := time.ParseDuration("1s")

	//	QuantizeFunction[fsEvent]
	_ = func(readyChan chan<- bool, queue []MyNiceEvent) {
		if len(queue) > 0 {
			readyChan <- true
		} else {
			time.Sleep(interval)
			readyChan <- false
		}
	}

}
