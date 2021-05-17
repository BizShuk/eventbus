package eventbus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetChannel() chan Event {
	args := m.Called()
	return args.Get(0).(chan Event)
}

func (m *MockService) GetEventType() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockService) Run() {
	_ = m.Called()
}

func (m *MockService) Stop() {
	_ = m.Called()
}

func TestEventDispatcher_Publish(t *testing.T) {
	mockEvent := &Event{Ctx: context.TODO(), Data: "TestMessage", EventType: "MockEventType"}
	mockEvent_EmptyEventType := &Event{Ctx: context.TODO(), Data: "TestMessage", EventType: ""}
	mockEvent_NilContext := &Event{Ctx: nil, Data: "TestMessage", EventType: "MockEventType"}

	type args struct {
		ctx    context.Context
		events []*Event
	}
	tests := []struct {
		name    string
		e       *EventDispatcher
		args    args
		wantErr bool
	}{
		{
			name: "TestEventDispatcher_Publish",
			e: &EventDispatcher{
				EventChan: make(map[string]EventService),
			},
			args:    args{ctx: context.TODO(), events: []*Event{mockEvent}},
			wantErr: false,
		},
		{
			name: "TestEventDispatcher_Publish_EventIsNil",
			e: &EventDispatcher{
				EventChan: make(map[string]EventService),
			},
			args:    args{ctx: context.TODO(), events: []*Event{nil}},
			wantErr: false,
		},
		{
			name: "TestEventDispatcher_Publish_EventTypeIsEmpty",
			e: &EventDispatcher{
				EventChan: make(map[string]EventService),
			},
			args:    args{ctx: context.TODO(), events: []*Event{mockEvent_EmptyEventType}},
			wantErr: false,
		},
		{
			name: "TestEventDispatcher_Publish_EventTypeIsEmpty",
			e: &EventDispatcher{
				EventChan: make(map[string]EventService),
			},
			args:    args{ctx: context.TODO(), events: []*Event{mockEvent_NilContext}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Publish(tt.args.ctx, tt.args.events...); (err != nil) != tt.wantErr {
				t.Errorf("EventDispatcher.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventDispatcher_dispatch(t *testing.T) {
	mockEvent := Event{Ctx: context.TODO(), Data: "TestMessage", EventType: "MockEventType"}
	mockService := &MockService{}
	mockService.On("GetChannel").Return(make(chan Event, 10))

	type args struct {
		event Event
	}
	tests := []struct {
		name string
		e    *EventDispatcher
		args args
	}{
		{
			name: "TestEventDispatcher_dispatch",
			e: &EventDispatcher{
				EventChan: map[string]EventService{
					"MockEventType": mockService,
				},
			},
			args: args{event: mockEvent},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.dispatch(tt.args.event)
		})
	}
}

func TestEventDispatcher_Registor(t *testing.T) {

	mockService := &MockService{}
	mockService.On("GetChannel").Return(make(chan Event, 10))
	mockService.On("GetEventType").Return("MockEventType")
	mockService.On("Run").Return()

	type args struct {
		svc EventService
	}
	tests := []struct {
		name    string
		e       *EventDispatcher
		args    args
		wantErr bool
	}{

		{
			name: "TestEventDispatcher_Registor_RegistorEventService",
			e: &EventDispatcher{
				EventChan: make(map[string]EventService),
			},
			args:    args{svc: mockService},
			wantErr: false,
		},
		{
			name: "TestEventDispatcher_Registor_HaveRegistoredEventService",
			e: &EventDispatcher{
				EventChan: map[string]EventService{
					"MockEventType": mockService,
				},
			},
			args:    args{svc: mockService},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Registor(tt.args.svc); (err != nil) != tt.wantErr {
				t.Errorf("EventDispatcher.Registor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventDispatcher_Unregistor(t *testing.T) {
	mockEventType := "MockEventType"
	mockService := &MockService{}
	mockService.On("GetChannel").Return(make(chan Event, 10))
	mockService.On("GetEventType").Return(mockEventType)
	mockService.On("Stop").Return()

	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		e    *EventDispatcher
		args args
	}{
		{
			name: "TestEventDispatcher_Unregistor",
			e: &EventDispatcher{
				EventChan: map[string]EventService{
					"MockEventType": mockService,
				},
			},
			args: args{eventType: mockEventType},
		},
		{
			name: "TestEventDispatcher_Unregistor_NoRegistedEventService",
			e: &EventDispatcher{
				EventChan: make(map[string]EventService),
			},
			args: args{eventType: mockEventType},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Unregistor(tt.args.eventType)
		})
	}
}
