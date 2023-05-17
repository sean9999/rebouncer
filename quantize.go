package rebouncer

// QuantizeFunction reads from the Queue, deciding when to emit(), and when to call itself again.
// QuantizeFunction is run any time `false` is written to readyChannel.
// Periodicity is achieved when QuantizeFunction itself writes to readyChannel
//
// A value of `false` sent to readyChannel triggers another run of QuantizeFunction.
// A value of `true` triggers emit()
type QuantizeFunction[T any] func(chan<- bool, []NiceEvent[T])
