package eventbus

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSEventService struct {
	ch         chan Event
	eventType  string
	client     *http.Client
	snsClient  SNSPublishAPI
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
	ctx, s.cancelFunc = context.WithCancel(context.TODO()) // This is for canceling go routine. TODO is enough

	go func(c context.Context, api SNSPublishAPI, topicArn string) {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-s.GetChannel():
				input := &sns.PublishInput{
					Message:  event.Data.(*string),
					TopicArn: &topicArn,
				}

				_, err := PublishMessage(event.Ctx, api, input) // _ (result), might be useful in some case
				if err != nil {
					log.Println("Got an error publishing the message:")
					log.Println(err)
					return
				}
			}
		}
	}(ctx, s.snsClient, s.topic)
}

func (s *SNSEventService) Stop() {
	if s.cancelFunc == nil {
		return
	}
	s.cancelFunc()
}

const DefaultSNSServiceEventType string = "SNSEvent"

func CreateSNSEventService(options ...SNSOption) (*SNSEventService, error) {
	svc := &SNSEventService{
		ch:        make(chan Event),
		eventType: DefaultSNSServiceEventType,
		client:    http.DefaultClient,
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(svc)
	}

	svc.createSNSClient()

	return svc, nil
}

func (svc *SNSEventService) createSNSClient() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithHTTPClient(svc.client))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	svc.snsClient = sns.NewFromConfig(cfg)

}

// Below is from AWS-go-sdk-v2 example

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
	if api == nil {
		return nil, errors.New("No SNS Client")
	}
	return api.Publish(c, input)
}
