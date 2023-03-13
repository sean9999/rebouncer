package main

import (
	"flag"
	"log"

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

	niceEvents := make(chan rebouncer.NiceEvent, rebouncer.BufferLength)

	if err := rebouncer.WatchDirectory(*watchDir, niceEvents); err != nil {
		log.Fatal(err)
	}

	for e := range niceEvents {
		log.Printf("Rebouncer: %s - %s", e.Event, e.File)
	}

}
