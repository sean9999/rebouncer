/*
Package rebouncer is a generic library that takes a noisy source of events
and produces a calmer, fitter, and happier source.

It has a well-defined set of lifecycle events that the user hooks into to get desired functionality.

The canonical example is a file-watcher that discards events involving temp files that IDEs might create. A file-watcher will also typically want a "rebounce" feature that prevents premature firing. Rebouncer provides a generic framework that can solve these problems by allowing the user to inject three types of user-defined functions: [Ingester], [Reducer], and [Quantizer].

# Components

These architectural components are involved in making Rebouncer work:

  - The NiceEvent is the atomic unit. It is a user-defined type. It is whatever you need it to be for your use case.
  - The [Ingester] produces events. When it's work is done, Rebouncer enters the [Draining] lifecycle state.
  - The [Reducer] is run every time after [Ingester] pushes an event to the [Queue]. It operates on all records in the queue and modifies the queue in its totality.
  - The [Queue] is a memory-safe slice of Events, waiting to be flushed to the consumer
  - The [Quantizer] runs at intervals of its choosing, deciding whether or not to flush to the consumer. It and [Reducer] take turns locking the [Queue], ensuring safety.

These mechanical components exist to enable the above:

  - an incomingEvents channel of type Event
  - a lifeCycle channel to keep track of lifecycle state.
  - a mutex lock to enable memory-safe operations against the [Queue].

# Behaviour

When [Ingester] completes, Rebouncer enters the [Draining] state.

You can receive events with [rebouncer.Subscribe], which returns a channel.

You can trigger the [Draining] state with [rebouncer.Interrupt].
*/
package rebouncer
