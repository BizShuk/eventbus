package eventbus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSEventService struct {
	ch         chan Event
	eventType  string
	client     *http.Client
	topic      string // AWS ARN
	cancelFunc context.CancelFunc
}

func (s *SNSEventService) GetChannel() chan Event {
	return s.ch
}

func (s *SNSEventService) GetEventType() string {
	return s.eventType
}

func (s *SNSEventService) Run() {
	if s.cancelFunc != nil {
		fmt.Println("cancel first")
		return
	}
	var ctx context.Context
	ctx, s.cancelFunc = context.WithCancel(context.TODO())
	// TODO: create SNS cilent by aws go sdk v2
	client := s.createSNSClient()

	go func() {
		topicArn := s.topic

		select {
		case <-ctx.Done():
		case event := <-s.GetChannel():

			input := &sns.PublishInput{
				Message:  event.Data.(*string),
				TopicArn: &topicArn,
			}

			_, err := PublishMessage(event.Ctx, client, input) // _ might be useful in some case
			if err != nil {
				fmt.Println("Got an error publishing the message:")
				fmt.Println(err)
				return
			}
		}
	}()
}

func (s *SNSEventService) Stop() {
	if s.cancelFunc == nil {
		return
	}
	s.cancelFunc()
}

const DefaultEventType string = "SNSEvent"

func CreateSNSEventService(options ...SNSOption) (*SNSEventService, error) {
	svc := &SNSEventService{
		ch:        make(chan Event),
		eventType: DefaultEventType,
		client:    http.DefaultClient,
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(svc)
	}

	return svc, nil
}

func (svc *SNSEventService) createSNSClient() *sns.Client {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithHTTPClient(svc.client))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)
	return client
}

// SNSPublishAPI defines the interface for the Publish function.
// We use this interface to test the function using a mocked service.
type SNSPublishAPI interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// PublishMessage publishes a message to an Amazon Simple Notification Service (Amazon SNS) topic
// Inputs:
//     c is the context of the method call, which includes the Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PublishOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to Publish
func PublishMessage(c context.Context, api SNSPublishAPI, input *sns.PublishInput) (*sns.PublishOutput, error) {
	return api.Publish(c, input)
}
