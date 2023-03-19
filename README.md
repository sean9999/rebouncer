# Rebouncer

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/rebouncer/graphs/commit-activity)

[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/sean9999/go/rebouncer)

[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/rebouncer)](https://goreportcard.com/report/github.com/sean9999/rebouncer)

[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/rebouncer.svg)](https://github.com/sean9999/rebouncer)

## A Powerful Debouncer for your Conjuring Needs

![Hand of Fish](/docs/hand.jpg)

Rebouncer is a generic library that takes a noisy source of events, and produces a cleaner source.

It employes a plugin architecture that can allow it to be used flexibly whenever the fan-out/fan-in concurrency pattern is needed.

The canonical case is a file-watcher that discards events involving temp files and other artefacts, providing its consumer with a clean, sane, and curated source of events. It is the engine behind [Fasthak](https://www.seanmacdonald.ca/posts/fasthak/).

For the canonical case, rebouncer is also available as a binary. It takes a directory as an argument, producing SSE events to stdout.

## Using it as a binary

This simplest case is accomplished like so:

```sh
$ go install github.com/sean9999/rebouncer/cmd
$ rebouncer -dir ~/projects/myapp/build
```

Which might stream to stdout something that looks like this:

<pre>
<samp>event: rebouncer/fs/output
data: {"file": "index.html", "operation": "modify"}

event: rebouncer/fs/output
data: {"file": "css/debug.css", "operation": "delete"}

event: rebouncer/fs/output
data: {"file": "css/mobile", "operation": "create"}

event: rebouncer/fs/output
data: {"file": "css/mobile", "operation": "modify"}</samp>
</pre>

## Using it as a library

You may want more flexibility than that. Rebouncer can be invoked as a library, allowing you to embed it in your application and giving you fine-grained control.

Rebouncer needs a few basic components, to be passed in. Let's go over them. Our examples will continue with the paradigm of building a file-watcher

### Injestor

An injestor is defined as runs in a go routine, and sends events of interest to Rebouncer, pushing them onto the Queue.

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

For more detailed docs, see [the docs](https://godoc.org/sean9999/go/rebouncer)
