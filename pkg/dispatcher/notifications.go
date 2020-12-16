package dispatcher

import (
	"context"
	"fmt"
	"github.com/we4tech/uampnotif/pkg/clients"
	"github.com/we4tech/uampnotif/pkg/configs"
	"github.com/we4tech/uampnotif/pkg/notifiers"
	"log"
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
	notificationCfg *notifiers.Config
	specCfg         map[string]*configs.Spec
	envVars         map[string]string
	events          chan DispatchEvent
	globalParams    map[string]string
	logger          *log.Logger
}

//
// Dispatch triggers events for all notifiers using separate go-routine.
//
func (n *notificationsDispatcher) Dispatch(ctx context.Context) error {
	n.logger.Println("Dispatching notifications")

	wg := sync.WaitGroup{}
	errs := make([]error, 0)
	errCh := make(chan error)

	go n.monitorEvents()
	go func() {
		for {
			select {
			case err := <-errCh:
				errs = append(errs, err)
			}
		}
	}()

	// TODO(HK): By default all requests are async.
	for _, notifier := range n.notificationCfg.Notifiers {
		n.logger.Printf("Dispatching %s\n", notifier.Id)
		wg.Add(1)

		go func(notifier notifiers.Notifier) {
			defer func() {
				n.logger.Printf("Dispatched %s\n", notifier.Id)
			}()

			defer wg.Done()

			if err := n.dispatchNotification(notifier); err != nil {
				errCh <- err
			}
		}(notifier)
	}

	wg.Wait()

	if len(errs) > 0 {
		return &dispatchError{Errors: errs}
	}

	return nil
}

func (n *notificationsDispatcher) isAsync(cfg *notifiers.Setting) bool {
	if cfg == nil {
		return n.notificationCfg.DefaultSettings.Async
	}

	return cfg.Async
}

func (n *notificationsDispatcher) SetMockClient(client clients.ClientImpl) {
	n.mockClient = client
}

func (n *notificationsDispatcher) dispatchNotification(notifier notifiers.Notifier) error {
	var (
		retries    = 0
		response   *clients.Response
		maxRetries = n.maxRetries(notifier.Settings)
	)

	go n.trigger(InTransit, notifier, retries, nil, nil)

restart:
	client, err := clients.NewHttpRequest(n.specCfg[notifier.Id], notifier.Params, n.envVars)
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
	notifier notifiers.Notifier,
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

func (n *notificationsDispatcher) maxRetries(settings *notifiers.Setting) int {
	if settings == nil {
		return n.notificationCfg.DefaultSettings.Retries
	}

	return settings.Retries
}

// Channel returns the DispatchEvent channel
func (n *notificationsDispatcher) Channel() chan DispatchEvent {
	return n.events
}

func (n *notificationsDispatcher) monitorEvents() {
	for {
		select {
		case event := <-n.Channel():
			n.logger.Printf("- %+v", event)
		}
	}
}

//
// NewNotificationDispatcher returns an implementation of *Dispatcher* service.
//
func NewNotificationDispatcher(
	logger *log.Logger,
	specs map[string]*configs.Spec,
	config *notifiers.Config,
	params, envVars map[string]string,
) Dispatcher {
	return &notificationsDispatcher{
		logger:          logger,
		specCfg:         specs,
		notificationCfg: config,
		globalParams:    params,
		envVars:         envVars,
		events:          make(chan DispatchEvent),
	}
}
