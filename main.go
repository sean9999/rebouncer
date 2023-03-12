package rebouncer

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/rjeczalik/notify"
)

// variables available to all functions
type Info struct {
	Dir string
}

// OriginalEvent can be anything
// Your mapper will turn it into a NiceEvent
type OriginalEvent interface{}

type OriginalEventChannel chan any

// NiceEvent is a notifyEvent, formatted nicely for our purposes
type NiceEvent struct {
	Tag   string
	Event string
	File  string
}

type Mapper func(chan any, chan NiceEvent) NiceEvent

func DefaultMapFunction(info Info, ei any) NiceEvent {
	abs, _ := filepath.Abs(info.Dir)
	specificEvent := ei.(notify.EventInfo)
	normalizedPath := strings.TrimPrefix(specificEvent.Path(), abs+"/")
	eventTag := "notify.EventInfo"
	return NiceEvent{
		Tag:   eventTag,
		File:  normalizedPath,
		Event: specificEvent.Event().String(),
	}
}

// Reducer takes an array of NiceEvents and returns an array of NiceEvents
// usually, filtering out duplicates or uninteresting events
type Reducer func(info Info, eventBatch []NiceEvent) []NiceEvent

func DefaultReduceFunction(info Info, inEvents []NiceEvent) []NiceEvent {
	outEvents := inEvents
	return outEvents
}

func WatchFolder(info Info, reduce Reducer) (chan NiceEvent, error) {
	var mapChannel chan NiceEvent
	var batchChannel chan []NiceEvent
	var rebounceChannel chan NiceEvent

	var fsEvents = make(chan notify.EventInfo)
	err := notify.Watch(info.Dir+"/...", fsEvents, notify.All)

	var thisBatch []NiceEvent

	go func(info Info) {
		select {
		case fsEvent := <-fsEvents:
			//	simply map original event to NiceEvent
			mapChannel <- DefaultMapFunction(info, fsEvent)
		case niceEvent := <-mapChannel:
			thisBatch = append(thisBatch, niceEvent)
		case <-time.After(3 * time.Second):
			transormedEvents := reduce(info, thisBatch)
			batchChannel <- transormedEvents
		case eventSlice := <-batchChannel:
			for _, e := range eventSlice {
				rebounceChannel <- e
			}
		}
	}(info)

	return rebounceChannel, err
}

func Setup(info Info, oChan chan any, mapFunc Mapper, reduceFunc Reducer) chan NiceEvent {

}

/*
func Rebounce(fsEvents chan NiceEvent, batchedEvents chan []NiceEvent, rebouncedEvents chan NiceEvent) {
	go Watch(".", fsEvents)
	go Quantize(fsEvents, batchedEvents)
	go Reduce(batchedEvents, rebouncedEvents)

	for ev := range rebouncedEvents {
		fmt.Println(ev)
	}

}
*/
