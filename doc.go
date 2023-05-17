/*
Package rebouncer is a generic library that takes a noisy source of events, and produces a calmer, fitter, and healthier source.

It employs a plugin architecture that can allow it to be used flexibly whenever the fan-out-fan-in concurrency pattern is needed.

The canonical case is a file-watcher that discards events involving temp files and other artefacts, providing it's consumer with a clean, sane, and curated source of events.

Example code will use that case.

Rebouncer can be thought of as a state machine with well-defined lifecycle states, and with dependency-injection for business logic. It makes liberal use of channels and go routines.
It's design ensures a healthy boundary between the generic mechanics that make Rebouncer work, and the business-logic that make it work for you and your use-case.

# Components

These general architectural components are involved:

  - The [NiceEvent] is the atomic unit. Most of the data flow through Rebouncer is with NiceEvents.
  - Queue is a buffer of NiceEvents, waiting to be flushed (emited) to the consumer
  - A readyChannel which is used to indicate when an emit() is approproate
  - NiceEvents are transferred from the Queue to a channel called outgoingEvents. This is what the consumer listens on.
  - A mutex to prevent loss of data through concurrency

# User-defined Functions

The general characteristics of user-defined functions in Rebouncer are:
  1. They only have access to the data they need (read-only when appropriate)
  2. Through the power of closures they can access scope that Rebouncer itself cannot
  3. They contain all business-logic necessary for the use-case.

User-defined functions are passed in at instantiation-time, and executed at well-defined times during lifecycle events.

The three types are [IngestFunction], [ReduceFunction], and [QuantizeFunction].

# Lifecycle Events

  - Ingest. Dirty events are transformed into NiceEvents and pushed to the Queue. An [IngestFunction] continually pumps data in.
  - Reduce. Operates on the entire Queue and replaces it wholesale. [ReduceFunction] runs any time the Queue is added to by [IngestFunction]
  - Quantize. Runs whenever it feels it's necessary, and decides to emit(), or not.
  - Egest. This is effectively just emit(), but it's important that a distinction be made between the act of emiting and the corresponding Lifecycle event, so we can reason about data loss.
*/
package rebouncer
