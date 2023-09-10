package rebouncer

import (
	"fmt"
	"time"
)

type fsEvent struct {
	File          string
	Operation     string
	TransactionId uint64
}

func ExampleQuantizer() {

	// type Quantizer[NICE any] func([]NICE) bool

	quantFunc := func(queue []fsEvent) bool {

		//	one second between runs
		time.Sleep(time.Second)

		//	return true if there is anything at all in the queue
		ok2flush := (len(queue) > 0)
		return ok2flush

	}

	fmt.Println(quantFunc([]fsEvent{}))
	// Output: false

}
