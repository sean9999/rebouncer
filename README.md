# Rebouncer

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/rebouncer/graphs/commit-activity)

[![Go Reference](https://pkg.go.dev/badge/github.com/sean9999/rebouncer.svg)](https://pkg.go.dev/github.com/sean9999/rebouncer)

[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/rebouncer)](https://goreportcard.com/report/github.com/sean9999/rebouncer)

[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/rebouncer.svg)](https://github.com/sean9999/rebouncer)

## A Powerful Debouncer for your Conjuring Needs

![Hand of Fish](/docs/hand.jpg)

Rebouncer is a generic library that takes a noisy source of events, and produces a cleaner source.

It employes a plugin architecture that can allow it to be used flexibly whenever the fan-out/fan-in concurrency pattern is needed.


## Using it as a library

Rebouncer needs a few basic to be passed in. Continuing the example a file-watcher, let's go over the basic architecture of these plugin lifecycle functions:

### Ingestor

An injestor is defined as runs in a go routine, and sends events of interest to Rebouncer, pushing them onto the Queue. It looks like this:

### Reducer

Reducer operates on the entire Queue, each time Injestor runs, modifiying, removing, or even adding events as needed.

### Quantizer

Runs in a go routine, keeping tabs on the Queue and telling Rebouncer when it's ready for a Flush().

The simplest case for use in a library is, again, the canonical case of a file-watcher, for which there is a convenience function:

```go
package main

import (
	"fmt"
	"github.com/sean9999/rebouncer"
)

//	watch ./build and emit every 1000 milliseconds
stateMachine := rebouncer.NewInotify("./build", 1000)

for niceEvent := range stateMachine.Subscribe() {
	fmt.Println(niceEvent.Dump())
}
```

Calling `rebouncer.NewInotify()` in this way is the equivilant of:

```go
//	rebecca is our singleton instance
stateMachine := rebouncer.New(rebouncer.Config{
	BufferSize: rebouncer.DefaultBufferSize,
	Quantizer:  rebouncer.DefaultInotifyQuantizer(1000),
	Reducer:    rebouncer.DefaultInotifyReduce,
	Injestor:   rebouncer.DefaultInotifyInjestor("./build", rebouncer.DefaultBufferSize),
})

```

`DefaultInotifyQuantizer()`, `DefaultInotifyReduce()`, and `DefaultInotifyInjestor()` are all themselves convenience functions that alleviate you from having to write your own respective Quantizer, Reducer, and Injestor.
