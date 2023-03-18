/*
Rebouncer is a generic library takes a noisy source of events, and produces a cleaner source.

It employs a plugin architecture that can allow it to be used flexibly whenever the fan-out-fan-in concurrency pattern is needed.

The canonical case is a file-watcher that discards events involving temp files and other artefacts, providing it's consumer with a clean, sane, and curated source of events.

For that case, rebouncer is also available as a binary, which takes a directory as an argument, producing SSE events to stdout.

To use the binary as a file-watcher:

	$ rebouncer -dir ./some/dir

To use it as a library, but again employing it as an inotify-backed filewatcher:

	stateMachine := rebouncer.NewInotify("./build", 1000)
	niceChannel := stateMachine.Subscribe()

Although Rebouncer provides convenience functions for common cases (such as file-watcher using inotify), an understanding of its basic architecture is necessary for more advanced uses.

  - an [Injestor] injests your source events, converting them into a format rebouncer can reason about, adding them to the queue
  - a [Reducer] operates on the entire queue of events, discarding, modifying, or even adding new ones at will
  - a [Quantizer] is initialized at startup and runs directly after the reducer, deciding where it's time to Emit()
  - an [Egestor] that formats the output. It is simply a function that takes a [NiceEvent] and returns an AnyEvent
*/
package rebouncer
