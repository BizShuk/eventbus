package eventbus

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWithEventType(t *testing.T) {
	mockSNSService := &SNSEventService{}
	mockEventType := "MockEventType"
	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestWithEventType",
			args: args{eventType: mockEventType},
			want: mockEventType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithEventType(tt.args.eventType)
			if got(mockSNSService); !reflect.DeepEqual(mockSNSService.eventType, tt.want) {
				t.Errorf("WithEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEventChan(t *testing.T) {

	mockSNSService := &SNSEventService{}
	ch := make(chan Event, 100)
	type args struct {
		ch chan Event
	}
	tests := []struct {
		name string
		args args
		want chan Event
	}{
		{
			name: "TestWithEventChan",
			args: args{ch: ch},
			want: ch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithEventChan(tt.args.ch)
			if got(mockSNSService); !reflect.DeepEqual(mockSNSService.ch, tt.want) {
				t.Errorf("WithEventChan() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestWithTopicArn(t *testing.T) {
	mockSNSService := &SNSEventService{}
	type args struct {
		arn string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestWithTopicArn",
			args: args{arn: "arn"},
			want: "arn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithTopicArn(tt.args.arn)
			if got(mockSNSService); !reflect.DeepEqual(mockSNSService.topic, tt.want) {
				t.Errorf("WithTopicArn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithHTTPClient(t *testing.T) {

	mockSNSService := &SNSEventService{}

	type args struct {
		client *http.Client
	}
	tests := []struct {
		name string
		args args
		want *http.Client
	}{
		{
			name: "TestWithHTTPClient",
			args: args{client: http.DefaultClient},
			want: http.DefaultClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithHTTPClient(tt.args.client)
			if got(mockSNSService); !reflect.DeepEqual(mockSNSService.client, tt.want) {
				t.Errorf("WithHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
