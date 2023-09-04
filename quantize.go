package rebouncer

// QuantizeFunction reads from the Queue, deciding when to emit(), and when to call itself again.
// QuantizeFunction is run any time `false` is written to readyChannel.
// Periodicity is achieved when QuantizeFunction itself writes to readyChannel
//
// A value of `false` sent to readyChannel triggers another run of QuantizeFunction.
// A value of `true` triggers emit()
type QuantizeFunction[NICE any] func([]NICE) bool

func (m *stateMachine[NICE]) quantize(fn QuantizeFunction[NICE]) {
	go func() {
		readyToEmit := fn(m.readQueue())
		if readyToEmit {
			m.lifeCycle <- Emiting
		} else {
			if m.GetLifeCycleState() == Draining && len(m.outgoingEvents) == 0 && len(m.queue) == 0 {
				m.lifeCycle <- Drained
			} else {
				m.lifeCycle <- Quantizing
			}
		}
	}()
}
