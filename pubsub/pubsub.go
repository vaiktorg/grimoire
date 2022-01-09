package pubsub

import (
	"log"
	"reflect"
	"sync"

	"github.com/vaiktorg/grimoire/errs"
)

//EventBus type of EventBus.
type EventBus struct {
	mutex    sync.Mutex
	handlers map[string][]reflect.Value
}

//NewEventBus instances an EventBus
func (e *EventBus) Init() {
	e.handlers = make(map[string][]reflect.Value)
}

// Notify sends data to all subscriptions.
func (e *EventBus) Notify(ev interface{}) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	evVal := reflect.ValueOf(ev)
	evTyp := evVal.Type().String()

	if !evVal.IsValid() {
		e.error("event is nil")
	}

	var hndlrs []reflect.Value
	if hs, ok := e.handlers[evTyp]; ok {
		hndlrs = hs
	} else {
		e.error("no registered handlers for this event type")
	}

	for _, hndlr := range hndlrs {
		hndlrEvType := hndlr.Type().In(0)
		if hndlrEvType.String() != evTyp {
			e.error("event type does not match subscribed handlers")
		}
		hndlr.Call([]reflect.Value{evVal})
	}

	return nil
}

//Subscribe adds a closure that handles the specific event.
func (e *EventBus) Subscribe(handlers ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, handler := range handlers {
		if handler != nil {
			evVal := reflect.ValueOf(handler)
			if evVal.Kind() != reflect.Func {
				panic(errs.Error("provided hanlder is not of type func"))
			}

			evTyp := evVal.Type().In(0).String()

			if _, OK := e.handlers[evTyp]; !OK {
				e.handlers[evTyp] = []reflect.Value{evVal}
				return
			}

			e.handlers[evTyp] = append(e.handlers[evTyp], evVal)
		}
	}
}

//Unsubscribe removes an event from our subscriptions.
func (e *EventBus) Unsubscribe(handlers ...interface{}) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, handler := range handlers {
		evType := reflect.TypeOf(handler).String()

		if _, OK := e.handlers[evType]; !OK {
			e.error("event type not registered")
		}

		delete(e.handlers, evType)
	}

	return nil
}
func (e *EventBus) Dispose() {
	e.handlers = nil
}
func (e *EventBus) error(err string) {
	log.Fatalf(err)
}
