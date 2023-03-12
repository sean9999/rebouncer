package rebouncer

// Quantizer takes a NiceEvent, adds it to a queue, and decides when to send off a batch
type Quantizer func(NiceEvent)

//var thisBatch []NiceEvent

// Quantize receives NiceEvents, batches them, and sends them to a channel taking BatchEvents
//
// @todo: seperate out a "Quantize" function that can be injected.
// @todo: make channels read-only or write-only as appropriate
/*
func Quantize(inchannel chan NiceEvent, outchannel chan []NiceEvent) {

	select {
	case e := <-inchannel:
		//	do some cleansing
		thisBatch = append(thisBatch, e)
	case <-time.After(3 * time.Second):
		outchannel <- thisBatch
		thisBatch = nil
	}
}
*/
