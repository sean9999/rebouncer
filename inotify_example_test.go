package rebouncer_test

import (
	"fmt"

	"github.com/sean9999/rebouncer"
)

// A convenience constructor suitable for the common case of watching a directory
func ExampleNewInotify() {

	//	watch ./build and emit every 1000 milliseconds
	stateMachine := rebouncer.NewInotify("./build", 1000)

	//	Subscribe() returns an active channel that NiceEvents are emitted to
	for niceEvent := range stateMachine.Subscribe() {
		fmt.Println(niceEvent.Dump())
	}

}
