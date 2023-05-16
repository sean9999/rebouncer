package rebouncer

// QuantizeFunction operates on a Queue and decides whether or not to flush it to the consumer
type QuantizeFunction[T any] func(chan<- bool, []NiceEvent[T])

// Quantizer runs in a go routine and sends true to readyChannel when it decides we're ready to emit()
// it has access to the entire batch in the Queue to help make this decision.
//type Quantizer[T any] func(chan bool, *[]NiceEvent[T])

func (m *machine[T]) quantize() {
	m.readyChannel <- false
}
