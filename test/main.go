package main

import (
	"context"
	"fmt"

	. "github.com/BizShuk/eventbus"
)

func main() {
	// Application config layer
	fmt.Println("abc")
	eb := &EventDispatcher{
		EventChan: make(map[string]EventService),
	}
	svc, err := CreateExampleEventService()
	if err != nil {
		fmt.Println(err)
	}

	err = eb.Registor(svc)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Service layer
	event := &Event{
		Data:      "abc",
		EventType: DefaultExampleEventType,
	}

	err = eb.Publish(context.Background(), event)
	if err != nil {
		fmt.Println(err)
		return
	}

	eb.Unregistor(DefaultExampleEventType)

}
