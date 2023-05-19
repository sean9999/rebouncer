/*
Package rebouncer is a generic library that takes a noisy source of events
and produces a calmer, fitter, and happier source.

It employs a plugin architecture that can allow it to be used flexibly.

The canonical case is a file-watcher that discards events involving temp files and other artefacts,
providing it's consumer with a clean, sane, and curated (not drinking too much) source of events.
This is useful with IDEs that tend to produce a lot of temp files or perform several operations on a file in quick succession.

Example code will use that case.

Rebouncer can be thought of as a state machine with well-defined lifecycle hooks
into which user-defined functions are injected.
It's design facilitates a healthy boundary between general mechanics and business logic.
It makes liberal use of channels and go routines.

# Components

These architectural components are involved in making Rebouncer work:

  - The [NiceEvent] is the atomic unit. Most of the data flow through Rebouncer is with NiceEvents. It wraps the original event (the noisy one) with a little metadata to help things along.
  - Queue is a buffer of NiceEvents, waiting to be flushed (emitted) to the consumer
  - A readyChannel which is used to indicate when an emit() is appropriate
  - NiceEvents are transferred from the Queue to a channel called outgoingEvents. This is what the consumer listens on.
  - A mutex to prevent loss of data through concurrency

The struct (machine) that implements [Behaviour] has no exported fields.
Behaviour has only one: [Behaviour.Subscribe] which is a channel of NiceEvents that Rebouncer's consumer listens on.

# User-defined Functions

The general characteristics of user-defined functions in Rebouncer are:
 1. They only have access to the data they need (read-only when appropriate)
 2. Through the power of closures they can access scope that Rebouncer itself cannot
 3. They contain all business-logic necessary for the use-case.
 4. They are triggered by Rebouncer at appropriate times.

User-defined functions are passed in at instantiation-time.

The three types of user-defined functions are [IngestFunction], [ReduceFunction], and [QuantizeFunction], corresponding to lifecycle events.

# Lifecycle Events

  - Ingest. Dirty events are transformed into NiceEvents and pushed to the Queue. An [IngestFunction] continually pumps data in.
  - Reduce. Operates on the entire Queue and replaces it wholesale. [ReduceFunction] runs any time the Queue is added to by [IngestFunction]
  - Quantize. Runs whenever it feels it's necessary, and decides to emit(), or not.
  - Egest. This is effectively just emit(), but it's important that a distinction be made between
    the act of emitting and the corresponding lifecycle event, so we can reason about data loss.
*/
package rebouncer
