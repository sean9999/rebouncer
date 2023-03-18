# Rebouncer

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/rebouncer/graphs/commit-activity)

[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/sean9999/go/rebouncer)

[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/rebouncer)](https://goreportcard.com/report/github.com/sean9999/rebouncer)

[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/rebouncer.svg)](https://github.com/sean9999/rebouncer)

## A Powerful Debouncer for your Conjuring

![Hand of Fish](/docs/hand.jpg)

Rebouncer is a generic library takes a noisy source of events, and produces a cleaner source.

It employes a plugin architecture that can allow it to be used flexibly whenever the fan-out/fan-in concurrency pattern is needed.

The canonical case is a file-watcher that discards events involving temp files and other artefacts, providing it's consumer with a clean, sane, and curated source of events. 

For that case, rebouncer is also available as a binary, which takes a directory as an argument, producing SSE events to stdout.

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

An injestor is defined as

```go
type Injestor[Subtype any] func() chan NiceEvent
```

This is where you listen on your original source of events, and for each event you get you transform it to a niceEvent. You pass back your NiceEvent channel, and it will continue to receive traffic.

A NiceEvent is our basic unit. It is a simple struct that we can rely on and reason about

```go
type NiceEvent[T any] struct {
    id unit64 // should be be unique
    Topic string // ex: "rebouncer/inotify/output
    Data T // it's up to you what to put here
}
```

Our Injestor might look like this:

```go
injestor := func() chan NiceEvent {
    niceEventChannel := make(chan, NiceEvent)
	
    //  these are the events we're injesting, but Rebouncer has no direct access to
    var fsEvents = make(chan notify.EventInfo, DefaultBufferSize)
	err := notify.Watch(dir+"/...", fsEvents, WatchMask)
	if err != nil {
		panic(err)
	}

    //  we're spawning a go routine, listening for notify events, transformatin them to NiceEvents, and pushing them to a channel which we return
	go func() {
		for fsEvent := range fsEvents {



			//m.Push(NotifyToNiceEvent(fsEvent, dir))
		}
	}() 
}
```