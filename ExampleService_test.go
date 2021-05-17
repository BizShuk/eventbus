package eventbus

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExampleEventService_GetChannel(t *testing.T) {

	ch := make(chan Event)
	tests := []struct {
		name string
		s    *ExampleEventService
		want chan Event
	}{
		{
			name: "TestExampleEventService_GetChannel",
			s:    &ExampleEventService{ch: ch},
			want: ch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExampleEventService.GetChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleEventService_GetEventType(t *testing.T) {
	tests := []struct {
		name string
		s    *ExampleEventService
		want string
	}{
		{
			name: "TestExampleEventService_GetEventType",
			s:    &ExampleEventService{eventType: DefaultExampleEventType},
			want: DefaultExampleEventType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetEventType(); got != tt.want {
				t.Errorf("ExampleEventService.GetEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleEventService_Run(t *testing.T) {
	ch := make(chan Event, 10)
	mockService := &ExampleEventService{
		ch: ch,
	}
	mockService_cancelFunNotNil := &ExampleEventService{
		ch:         ch,
		cancelFunc: func() {},
	}
	tests := []struct {
		name string
		s    *ExampleEventService
	}{
		{
			name: "TestExampleEventService_Run",
			s:    mockService,
		},
		{
			name: "TestExampleEventService_Run_cancelFunNotNil",
			s:    mockService_cancelFunNotNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Run()
			ch <- Event{
				Data: "test message",
			}
			time.Sleep(100 * time.Millisecond)
			tt.s.Stop()
		})
	}
}

func TestExampleEventService_Stop(t *testing.T) {
	tests := []struct {
		name string
		s    *ExampleEventService
	}{
		{
			name: "TestExampleEventService_Stop_cancelFuncIsNil",
			s:    &ExampleEventService{},
		},
		{
			name: "TestExampleEventService_Stop_cancelFuncIsNotNil",
			s:    &ExampleEventService{cancelFunc: func() {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Stop()
		})
	}
}

func TestCreateExampleEventService(t *testing.T) {
	tests := []struct {
		name    string
		want    *ExampleEventService
		wantErr bool
	}{
		{
			name:    "TestCreateExampleEventService",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateExampleEventService()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateExampleEventService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
