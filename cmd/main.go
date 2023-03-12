package main

import (
	"flag"

	"github.com/rjeczalik/notify"
	"github.com/sean9999/rebouncer"
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

	fsEvents := make(chan any)
	rootDirectory := *watchDir + "/..."
	err := notify.Watch(rootDirectory, fsEvents, notify.All)
	if err != nil {
		panic(err)
	}

	info := rebouncer.Info{Dir: rootDirectory}

	var mapper rebouncer.Mapper

	rebouncedEvents := rebouncer.Setup(info, fsEvents, rebouncer.DefaultMapFunction, rebouncer.DefaultReduceFunction)

}
