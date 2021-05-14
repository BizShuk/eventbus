package eventbus

import (
	"context"
	"fmt"
)

type ExampleEventService struct {
	ch         chan Event
	eventType  string
	cancelFunc context.CancelFunc
}

func (s *ExampleEventService) GetChannel() chan Event {
	return s.ch
}

func (s *ExampleEventService) GetEventType() string {
	return s.eventType
}

func (s *ExampleEventService) Run() {
	if s.cancelFunc != nil {
		fmt.Println("cancel first")
		return
	}
	var ctx context.Context
	ctx, s.cancelFunc = context.WithCancel(context.TODO()) // This is for canceling go routine. TODO is enough

	go func(c context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-s.GetChannel():
				fmt.Println(event)
			}
		}
	}(ctx)
}

func (s *ExampleEventService) Stop() {
	if s.cancelFunc == nil {
		return
	}
	s.cancelFunc()
}

const DefaultExampleEventType string = "ExampleEvent"

func CreateExampleEventService() (*ExampleEventService, error) {
	svc := &ExampleEventService{
		ch:        make(chan Event),
		eventType: DefaultExampleEventType,
	}

	return svc, nil
}
