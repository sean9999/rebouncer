package main

import (
	"flag"
	"log"

	rebouncer "github.com/sean9999/rebouncer/pkg"
)

var (
	watchDir *string
)

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	flag.Parse()
}

func main() {

	niceEvents := make(chan rebouncer.NiceEvent, rebouncer.BUFFER_LENGTH)

	//	start watcher
	if err := rebouncer.WatchRecursively(*watchDir, niceEvents); err != nil {
		log.Fatal(err)
	}

	/*
		for {
			select {
			case e := <-niceEvents:
				log.Printf("%s - %s", e.Event, e.File)
			}
		}
	*/

	for e := range niceEvents {
		log.Printf("MAIN NICEEVENT: %s - %s", e.Event, e.File)
	}

}
