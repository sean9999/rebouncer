package main

import (
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
	flushPeriod = flag.Int("period", 30000, "how often (in milliseconds) to flush events")
	flag.Parse()
}

func main() {

	//	instantiate
	//rebecca := rebouncer.NewInotify(*watchDir, *flushPeriod)

	rebecca := rebouncer.New(rebouncer.Config{
		BufferSize: 1024,
		Quantizer:  rebouncer.DefaultInotifyQuantizer(*flushPeriod),
		Reducer:    rebouncer.DefaultInotifyReducer,
	})
	go rebecca.WatchDir(*watchDir)

	//	here is the channel we can listen on
	outgoingEvents := rebecca.Subscribe()

	//	for example
	for e := range outgoingEvents {
		fmt.Println(e.Dump())
	}

}
