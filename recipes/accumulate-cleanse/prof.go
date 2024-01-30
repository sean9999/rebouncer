package main

import (
	"fmt"
	"os"
	"runtime/pprof"
)

func init() {
	f, err := os.Create("myprogram.prof")
	if err != nil {

		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
}
