package eventbus

import "context"

// prepare in service and pass in event publisher method
// ctx will be assigned if ctx is nil
type Event struct {
	Ctx       context.Context
	Data      interface{}
	EventType string
}
