package eventbus

import (
	"net/http"
)

type SNSOption func(svc *SNSEventService)

func WithEventType(eventType string) SNSOption {
	return func(svc *SNSEventService) {
		svc.eventType = eventType
	}
}

func WithEventChan(ch chan Event) SNSOption {
	return func(svc *SNSEventService) {
		svc.ch = ch
	}
}

func WithTopicArn(arn string) SNSOption {
	return func(svc *SNSEventService) {
		svc.topic = arn
	}
}

// XRay client can inject here
func WithHTTPClient(client *http.Client) SNSOption {
	return func(svc *SNSEventService) {
		svc.client = client
	}
}
