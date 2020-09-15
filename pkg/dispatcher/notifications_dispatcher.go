package dispatcher

import "github.com/we4tech/uampnotif/pkg/notifications"

type notificationsDispatcher struct {
	config *notifications.Config
}

func (n *notificationsDispatcher) Dispatch(events chan DispatchEvent) {
	panic("implement me")
}

//
// NewNotificationDispatcher returns an implementation of *Dispatcher* service.
//
func NewNotificationDispatcher(config *notifications.Config) Dispatcher {
	d := &notificationsDispatcher{
		config: config,
	}

	return d
}
