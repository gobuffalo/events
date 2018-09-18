package events

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

const (
	// AppStart is emitted when buffalo.App#Serve is called
	AppStart = "app:start"
	// AppStop is emitted when buffalo.App#Stop is called
	AppStop = "app:stop"
	// WorkerStart is emitted when buffalo.App#Serve is called and workers are started
	WorkerStart = "worker:start"
	// WorkerStop is emitted when buffalo.App#Stop is called and workers are stopped
	WorkerStop = "worker:stop"
	// RouteStarted is emitted when a requested route is being processed
	RouteStarted = "route:started"
	// RouteFinished is emitted when a requested route is completed
	RouteFinished = "route:finished"
	// ErrRoute is emitted when there is a problem handling processing a route
	ErrRoute = "err:route"
	// ErrGeneral is emitted for general errors
	ErrGeneral = "err:general"
	// ErrPanic is emitted when a panic is recovered
	ErrPanic = "err:panic"
	// ErrAppStart is emitted when an error occurs calling buffalo.App#Serve
	ErrAppStart = "err:app:start"
	// ErrAppStop is emitted when an error occurs calling buffalo.App#Stop
	ErrAppStop = "err:app:stop"
	// ErrWorkerStart is emitted when an error occurs when starting workers
	ErrWorkerStart = "err:worker:start"
	// ErrWorkerStop is emitted when an error occurs when stopping workers
	ErrWorkerStop = "err:worker:stop"
)

// Emit an event to all listeners
func Emit(e Event) error {
	return boss.Emit(e)
}

func EmitPayload(kind string, payload interface{}) error {
	return EmitError(kind, nil, payload)
}

func EmitError(kind string, err error, payload interface{}) error {
	e := Event{
		Kind:    kind,
		Payload: payload,
		Error:   err,
	}
	return Emit(e)
}

// NamedListen for events. Name is the name of the
// listener NOT the events you want to listen for,
// so something like "my-listener", "kafka-listener", etc...
func NamedListen(name string, l Listener) (DeleteFn, error) {
	return boss.Listen(name, l)
}

// Listen for events.
func Listen(l Listener) (DeleteFn, error) {
	_, file, line, _ := runtime.Caller(1)
	return NamedListen(fmt.Sprintf("%s:%d", file, line), l)
}

type listable interface {
	List() ([]string, error)
}

// List all listeners
func List() ([]string, error) {
	if l, ok := boss.(listable); ok {
		return l.List()
	}
	return []string{}, errors.Errorf("manager %T does not implemented listable", boss)
}
