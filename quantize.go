package rebouncer

// QuantizeFunction operates on a (read-only) Queue [machine.readyChannel] is written to with a value of false
// It _writes_ to readyChannel when it wants to run itself again or wants to signal that it's time to emit()
//
// Do NOT use a time.Ticker inside your Quantizer
// Periodicity is achieved by waiting and then writing `false` to readyChannel
type QuantizeFunction[T any] func(chan<- bool, []NiceEvent[T])
