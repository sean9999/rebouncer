package main

import (
	"flag"
	"fmt"

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

	rebel := rebouncer.New(rebouncer.Config{
		BufferSize: 1024,
	})

	niceEvents := rebel.Subscribe()
	go rebel.WatchDir(*watchDir)

	for e := range niceEvents {
		fmt.Println(e.Dump())
	}

}
