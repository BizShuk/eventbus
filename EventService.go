package eventbus

type EventService interface {
	GetChannel() chan Event
	GetEventType() string
	Run()
	Stop()
}
