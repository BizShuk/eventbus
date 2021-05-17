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

func TestWithMaxChannelBuffer(t *testing.T) {

	mockSNSService := &SNSEventService{}
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TestWithMaxChannelBuffer",
			args: args{num: 100},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithMaxChannelBuffer(tt.args.num)
			if got(mockSNSService); !reflect.DeepEqual(cap(mockSNSService.ch), tt.want) {
				t.Errorf("WithMaxChannelBuffer() = %v, want %v", got, tt.want)
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
