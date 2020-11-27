package dispatcher

import (
	"context"
	"github.com/we4tech/uampnotif/pkg/clients"
)

//
// DispatchState provides a dispatch state
//
type DispatchState string

const (
	InTransit = DispatchState("in-transit")
	Retrying  = DispatchState("retrying")
	Error     = DispatchState("error")
	Success   = DispatchState("success")
)

//
// DispatchEvent provides a struct for storing dispatching events.
//
type DispatchEvent struct {
	State        DispatchState
	NotifierId   string
	NotifierDesc string
	Retries      int
	Error        error
	Response     *clients.Response
}

//
// Dispatcher provides an interface for implementing notifier's dispatcher.
//
type Dispatcher interface {
	Dispatch(ctx context.Context) error
	Channel() chan DispatchEvent
	SetMockClient(impl clients.ClientImpl)
}
