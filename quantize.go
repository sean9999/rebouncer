package rebouncer

// QuantizeFunction reads from the Queue, deciding when to emit(), and when to call itself again.
// QuantizeFunction is run any time `false` is written to readyChannel.
// Periodicity is achieved when QuantizeFunction itself writes to readyChannel
//
// A value of `false` sent to readyChannel triggers another run of QuantizeFunction.
// A value of `true` triggers emit()
type QuantizeFunction[NICE any] func(chan<- bool, []NICE)

func (m *stateMachine[NICE]) quantize(fn QuantizeFunction[NICE]) {
	if m.lifeCycleState < Draining {
		m.lifeCycleState = Quantizing
	}
	if m.lifeCycleState < Drained {
		go fn(m.readyChannel, m.readQueue())
	}
}
