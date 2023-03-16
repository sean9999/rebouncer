package main

import (
	"flag"
	"fmt"
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

	rebel := rebouncer.New(rebouncer.Config{
		BufferSize: 1024,
	})

	fmt.Println("*** Rebouncer ***")
	fmt.Printf("%+v\n", rebel.Info())

	niceEvents := rebel.Subscribe()

	if err := rebouncer.WatchDirectory(*watchDir, niceEvents); err != nil {
		log.Fatal(err)
	}

	for e := range niceEvents {
		fmt.Println(e.Dump())
	}

}
