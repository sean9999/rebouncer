package rebouncer

// Quantizer reads from the Queue, deciding when to emit(), and when to call itself again.
// Quantizer is run any time `false` is written to readyChannel.
// Periodicity is achieved when Quantizer itself writes to readyChannel
//
// A value of `false` sent to readyChannel triggers another run of Quantizer.
// A value of `true` triggers emit()
type Quantizer[NICE any] func([]NICE) bool

func (m *rebounceMachine[NICE]) quantize(fn Quantizer[NICE]) {
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
