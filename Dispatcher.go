package eventbus

import (
	"context"
	"errors"
)

type Dispatcher interface {
	Publish(ctx context.Context, events ...*Event) error
	Registor(EventService) error
	UnRegistor(eventType string)
}

type EventDispatcher struct {
	EventChan map[string]EventService
}

// Here is gin context come in for X-ray or others
func (e *EventDispatcher) Publish(ctx context.Context, events ...*Event) error {
	var validEvents []Event

	for _, event := range events {
		if event == nil || event.EventType == "" {
			continue
		}
		if event.Ctx == nil { // inject contetx from service if not exist
			event.Ctx = ctx
		}
		validEvents = append(validEvents, *event)
	}

	for _, event := range validEvents {
		e.dispatch(event)
	}

	return nil
}

func (e *EventDispatcher) dispatch(event Event) {
	if svc, exist := e.EventChan[event.EventType]; exist {
		svc.GetChannel() <- event
	}
}

func (e *EventDispatcher) Registor(svc EventService) error {
	if _, exist := e.EventChan[svc.GetEventType()]; exist {
		return errors.New("EventType has been registered")
	}

	e.EventChan[svc.GetEventType()] = svc
	svc.Run()
	return nil
}

func (e *EventDispatcher) Unregistor(eventType string) {
	if svc, exist := e.EventChan[eventType]; exist {
		svc.Stop()
	}
	delete(e.EventChan, eventType)
}
