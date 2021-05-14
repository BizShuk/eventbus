package eventbus

type EventService interface {
	GetChannel() chan interface{}
	GetEventType() string
	Run()
	Stop()
}
