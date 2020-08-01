package coordinator

import (
	"time"
)

type EventAggregator struct {
	listeners map[string][]func(EventData)
}

type EventData struct {
	Name      string
	Value     float64
	TimeStamp time.Time
}

//NewEventAggregator intializes a new EventAggregator
func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(EventData)),
	}

	return &ea
}

//AddListener will be used by the consumer to register callback functions
//with the aggregator
func (ea *EventAggregator) AddListener(eventName string, f func(EventData)) {
	ea.listeners[eventName] = append(ea.listeners[eventName], f)
}

//PublishEvent publishes the EventData
func (ea *EventAggregator) PublishEvent(eventName string, eventData EventData) {
	if ea.listeners[eventName] != nil {
		//call all the cb func for that event
		for _, cb := range ea.listeners[eventName] {
			cb(eventData) //pass a copy for only read access for consumers
		}

	}
}
