package eventbus

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
)

func TestSNSEventService_GetChannel(t *testing.T) {
	ch := make(chan Event, 10)
	svc := &SNSEventService{
		ch: ch,
	}
	tests := []struct {
		name string
		s    *SNSEventService
		want chan Event
	}{
		{
			name: "TestSNSEventService_GetChannel",
			s:    svc,
			want: ch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SNSEventService.GetChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSNSEventService_GetEventType(t *testing.T) {
	svc := &SNSEventService{
		eventType: DefaultSNSServiceEventType,
	}
	tests := []struct {
		name string
		s    *SNSEventService
		want string
	}{
		{
			name: "TestSNSEventService_GetChannel",
			s:    svc,
			want: DefaultSNSServiceEventType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetEventType(); got != tt.want {
				t.Errorf("SNSEventService.GetEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSNSEventService_Run(t *testing.T) {
	snsClient := &mockSNSClient{}
	snsClient_ErrorOut := &mockSNSClient_ErrorOut{}
	ch := make(chan Event, 10)

	tests := []struct {
		name string
		s    *SNSEventService
	}{
		{
			name: "TestSNSEventService_Run_cancelFunNotNil",
			s: &SNSEventService{
				ch:         ch,
				cancelFunc: func() {},
				topic:      "arn",
			},
		},
		{
			name: "TestSNSEventService_Run_snsClientErrorOut",
			s: &SNSEventService{
				ch:        ch,
				topic:     "arn",
				snsClient: snsClient_ErrorOut,
			},
		},
		{
			name: "TestSNSEventService_Run",
			s: &SNSEventService{
				ch:        ch,
				topic:     "arn",
				snsClient: snsClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Run()
			testMsg := "test message"
			ch <- Event{
				Data: &testMsg,
			}
			time.Sleep(100 * time.Millisecond)
			tt.s.Stop()
		})
	}
}

func TestSNSEventService_Stop(t *testing.T) {
	tests := []struct {
		name string
		s    *SNSEventService
	}{
		{
			name: "TestSNSEventService_Stop",
			s:    &SNSEventService{},
		},
		{
			name: "TestSNSEventService_Stop_cancelFuncIsNotNil",
			s:    &SNSEventService{cancelFunc: func() {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Stop()
		})
	}
}

func TestCreateSNSEventService(t *testing.T) {
	type args struct {
		options []SNSOption
	}
	tests := []struct {
		name    string
		args    args
		want    *SNSEventService
		wantErr bool
	}{
		{
			name:    "TestCreateSNSEventService",
			args:    args{options: []SNSOption{nil, WithHTTPClient(http.DefaultClient)}},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSNSEventService(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSNSEventService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestSNSEventService_createSNSClient(t *testing.T) {
	tests := []struct {
		name string
		svc  *SNSEventService
		want *sns.Client
	}{
		{
			name: "TestSNSEventService_createSNSClient",
			svc:  &SNSEventService{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.svc.createSNSClient()
			assert.NotNil(t, tt.svc.snsClient)
		})
	}
}

type mockSNSClient struct{}

func (m *mockSNSClient) Publish(ctx context.Context, params *sns.PublishInput, optFuns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return &sns.PublishOutput{}, nil
}

type mockSNSClient_ErrorOut struct{}

func (m *mockSNSClient_ErrorOut) Publish(ctx context.Context, params *sns.PublishInput, optFuns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return nil, errors.New("something wrong during publish")
}

func TestPublishMessage(t *testing.T) {
	// snsClient := &mockSNSClient{} // slice spread operator cause some problem in tesify framework
	// snsClient.On("Publish", snsClient, &sns.PublishInput{}, func(*sns.Options) {}).Return(&sns.PublishOutput{}, nil)

	type args struct {
		c     context.Context
		api   SNSPublishAPI
		input *sns.PublishInput
	}
	tests := []struct {
		name    string
		args    args
		want    *sns.PublishOutput
		wantErr bool
	}{
		{
			name:    "TestPublishMessage",
			args:    args{c: context.TODO(), api: &mockSNSClient{}, input: &sns.PublishInput{}},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "TestPublishMessage_apiIsNil",
			args:    args{c: context.TODO(), api: nil, input: &sns.PublishInput{}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublishMessage(tt.args.c, tt.args.api, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublishMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			assert.NotNil(t, got)
		})
	}
}
