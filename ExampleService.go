package eventbus

import (
	"context"
	"log"
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
		log.Println("cancel first")
		return
	}
	var ctx context.Context
	ctx, s.cancelFunc = context.WithCancel(context.TODO()) // This is for canceling go routine. TODO is enough

	go func(c context.Context) {
		for {
			select {
			case <-ctx.Done():
				log.Println("ExampleService goroutine exits")
				return
			case event := <-s.GetChannel():
				log.Printf("event received from channel %v", event)
			}
		}
	}(ctx)
}

func (s *ExampleEventService) Stop() {
	if s.cancelFunc == nil {
		return
	}
	s.cancelFunc()

	log.Println("Called CancelFunc")
}

const DefaultExampleEventType string = "ExampleEvent"

func CreateExampleEventService() (*ExampleEventService, error) {
	svc := &ExampleEventService{
		ch:        make(chan Event),
		eventType: DefaultExampleEventType,
	}

	return svc, nil
}
