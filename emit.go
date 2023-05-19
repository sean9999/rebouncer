package rebouncer

// The Emit lifecycle event happens when a value is sent to readyChannel

type EmitFunction[NICE any, BEAUTIFUL any] func(NICE) BEAUTIFUL

func (m *stateMachine[NAUGHTY, NICE, BEAUTIFUL]) emit(fn EmitFunction[NICE, BEAUTIFUL]) {
	for _, niceEvent := range m.queue {
		m.outgoingEvents <- fn(niceEvent)
		//	@todo: provide a way to drain events one by one
	}
	m.writeQueue([]NICE{})
	m.readyChannel <- false // restart quantizer
}
