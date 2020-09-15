package dispatcher

//
// DispatchState provides a dispatch state
//
type DispatchState string

const (
	InTransit = DispatchState("in-transit")
	Retrying  = DispatchState("retrying")
	Error     = DispatchState("error")
)

//
// DispatchEvent provides a struct for storing dispatching events.
//
type DispatchEvent struct {
	State      DispatchState
	NotifierId string
	Retries    int
	Error      error
}

//
// Dispatcher provides an interface for implementing notifier's dispatcher.
//
type Dispatcher interface {
	Dispatch(chan DispatchEvent)
}
