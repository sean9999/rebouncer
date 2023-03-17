package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sean9999/rebouncer"
)

var watchDir *string

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	flag.Parse()
}

func main() {

	//	instantiate
	rebel := rebouncer.New(rebouncer.Config{
		BufferSize: 1024,
	})

	//	start the watcher
	go rebel.WatchDir(*watchDir)

	//	start a ticker

	tick := time.NewTicker(time.Hour)

	go func() {
		for t := range tick.C {
			fmt.Println("tick", t)
		}
	}()

	//	here is the channel we can listen on
	outgoingEvents := rebel.Subscribe()

	for e := range outgoingEvents {
		tick.Reset(3 * time.Second)
		fmt.Println(e.Dump())
	}

}
