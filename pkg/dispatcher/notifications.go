package dispatcher

import (
	"context"
	"fmt"
	"github.com/we4tech/uampnotif/pkg/clients"
	"github.com/we4tech/uampnotif/pkg/integrations"
	"github.com/we4tech/uampnotif/pkg/notifications"
	"sync"
)

type dispatchError struct {
	Errors []error
}

func (d *dispatchError) Error() string {
	return fmt.Sprintf("DispatchError: %+v\n", d.Errors)
}

type notificationsDispatcher struct {
	mockClient      clients.ClientImpl
	notificationCfg *notifications.Config
	specCfg         map[string]*integrations.Spec
	params          map[string]string
	envVars         map[string]string
	events          chan DispatchEvent
}

func (n *notificationsDispatcher) Dispatch(ctx context.Context) error {
	wg := sync.WaitGroup{}
	errCh := make(chan error)

	// TODO(HK): By default all requests are async.
	for _, notifier := range n.notificationCfg.Notifiers {
		go func(notifier notifications.Notifier) {
			wg.Add(1)
			defer wg.Done()

			if err := n.dispatchNotification(notifier); err != nil {
				errCh <- err
			}
		}(notifier)
	}

	wg.Wait()

	if len(errCh) == 0 {
		return nil
	}

	err := &dispatchError{Errors: make([]error, len(errCh))}
	for i := 0; i < len(errCh); i++ {
		err.Errors[i] = <-errCh
	}

	return err
}

func (n *notificationsDispatcher) isAsync(cfg *notifications.Setting) bool {
	if cfg == nil {
		return n.notificationCfg.DefaultSettings.Async
	}

	return cfg.Async
}

func (n *notificationsDispatcher) SetMockClient(client clients.ClientImpl) {
	n.mockClient = client
}

func (n *notificationsDispatcher) dispatchNotification(notifier notifications.Notifier) error {
	var (
		retries    = 0
		response   *clients.Response
		maxRetries = n.maxRetries(notifier.Settings)
	)

	go n.trigger(InTransit, notifier, retries, nil, nil)

restart:
	client, err := clients.NewHttpRequest(n.specCfg[notifier.Id], n.params, n.envVars)
	if err != nil {
		goto errorHandler
	}

	if n.mockClient != nil {
		client.SetClient(n.mockClient)
	}

	response, err = client.SendRequest()
	if err != nil {
		goto errorHandler
	}

	if response.IsOK() {
		go n.trigger(Success, notifier, retries, response, nil)
		return nil
	}

	if retries < maxRetries {
		go n.trigger(Retrying, notifier, retries, response, nil)

		retries++

		goto restart
	}

errorHandler:
	go n.trigger(Error, notifier, retries, response, err)

	return err
}

func (n *notificationsDispatcher) trigger(
	state DispatchState,
	notifier notifications.Notifier,
	retries int,
	response *clients.Response,
	err error,
) {
	n.events <- DispatchEvent{
		State:        state,
		NotifierId:   notifier.Id,
		NotifierDesc: notifier.Desc,
		Retries:      retries,
		Response:     response,
		Error:        err,
	}
}

func (n *notificationsDispatcher) maxRetries(settings *notifications.Setting) int {
	if settings == nil {
		return n.notificationCfg.DefaultSettings.Retries
	}

	return settings.Retries
}

// Channel returns the DispatchEvent channel
func (n *notificationsDispatcher) Channel() chan DispatchEvent {
	return n.events
}

//
// NewNotificationDispatcher returns an implementation of *Dispatcher* service.
//
func NewNotificationDispatcher(
	specs map[string]*integrations.Spec,
	config *notifications.Config,
	params, envVars map[string]string,
) Dispatcher {
	return &notificationsDispatcher{
		specCfg:         specs,
		notificationCfg: config,
		params:          params,
		envVars:         envVars,
		events:          make(chan DispatchEvent),
	}
}
