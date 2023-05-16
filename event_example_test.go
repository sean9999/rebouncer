package rebouncer_test

import (
	"fmt"
	"time"

	"github.com/sean9999/rebouncer"
)

// create a new NiceEvent
func ExampleNewNiceEvent() {
	e := rebouncer.NewNiceEvent("rebouncer/myplugin/incoming")
	fmt.Println(e)
}

// Using NewNiceEvent in an ingestor
func Exampleingestor() {

	var returnChannel = make(chan rebouncer.NiceEvent)

	ingestFunc := func() chan rebouncer.NiceEvent {

		nTicks := 0
		ticker := time.NewTicker(500 * time.Millisecond)

		go func() {
			for t := range ticker.C {
				e := rebouncer.NewNiceEvent("rebouncer/ticker/incoming")
				nTicks++
				fmt.Println(t, nTicks)
				returnChannel <- e
				if nTicks >= 10 {
					ticker.Stop()
					break
				}
			}
		}()

		return returnChannel

	}

	config := rebouncer.Config{
		ingestor: ingestFunc,
	}
	stateMachine := rebouncer.New(config)

	//	Subscribe() returns an active channel that NiceEvents are emitted to
	for niceEvent := range stateMachine.Subscribe() {
		fmt.Println(niceEvent.Dump())
	}

}
