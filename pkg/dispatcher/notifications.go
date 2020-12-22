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
	mutex           sync.Mutex
	done            chan struct{}
	globalParams    map[string]string
	logger          *log.Logger
}

//
// Done returns an error channel, data is only arrive whenever the dipatching process is done.
//
func (n *notificationsDispatcher) Done() chan struct{} {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.done == nil {
		n.done = make(chan struct{})
	}

	return n.done
}

//
// Dispatch triggers events for all notifiers using separate go-routine.
//
// TODO(HK): Add support for SYNC requests.
func (n *notificationsDispatcher) Dispatch(_ context.Context) error {
	n.logger.Println("Dispatching notifications")

	wg := &sync.WaitGroup{}
	errCh := make(chan error)

	wg.Add(len(n.notificationCfg.Notifiers))

	go n.dispatchInAsync(errCh, wg)

	wg.Wait()

	close(errCh)

	errs := n.accumulateErrors(errCh)

	if n.events != nil {
		close(n.events)
	}

	if n.done != nil {
		close(n.done)
	}

	if len(errs) > 0 {
		n.logger.Println("Failed to dispatch all notifications")

		return &dispatchError{Errors: errs}
	} else {
		n.logger.Println("Successfully dispatched all notifications")
	}

	return nil
}

func (n *notificationsDispatcher) SetMockClient(client clients.ClientImpl) {
	n.mockClient = client
}

// TODO(HK): Use context
func (n *notificationsDispatcher) dispatchNotification(
	_ context.Context,
	notifier notifiers.Notifier,
) error {
	var (
		retries    = 0
		response   *clients.Response
		maxRetries = n.maxRetries(notifier.Settings)
	)

	n.trigger(InTransit, notifier, retries, nil, nil)

restart:
	client, err := clients.NewHttpRequest(n.specCfg[notifier.Id], notifier.Params, n.envVars)
	if err != nil {
		goto errorHandler
	}

	if n.mockClient != nil {
		client.SetClient(n.mockClient)
	}

	// TODO(HK): Handle template error vs network error
	response, err = client.SendRequest()
	if err != nil {
		goto errorHandler
	}

	if response.IsOK() {
		n.trigger(Success, notifier, retries, response, nil)
		return nil
	}

	if retries < maxRetries {
		n.trigger(Retrying, notifier, retries, response, nil)

		retries++

		goto restart
	}

errorHandler:
	n.trigger(Error, notifier, retries, response, err)

	return err
}

func (n *notificationsDispatcher) trigger(
	state DispatchState,
	notifier notifiers.Notifier,
	retries int,
	response *clients.Response,
	err error,
) {
	if n.events == nil {
		return
	}

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

// Channel returns a lazily initiated DispatchEvent channel.
func (n *notificationsDispatcher) Events() chan DispatchEvent {
	if n.events != nil {
		return n.events
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.events == nil {
		n.events = make(chan DispatchEvent)
	}

	return n.events
}

func (n *notificationsDispatcher) dispatchInAsync(errCh chan error, wg *sync.WaitGroup) {
	ctx := context.Background()

	for _, notifier := range n.notificationCfg.Notifiers {
		n.logger.Printf("Dispatching %s\n", notifier.Id)

		go func(notifier notifiers.Notifier) {
			defer wg.Done()

			if err := n.dispatchNotification(ctx, notifier); err != nil {
				errCh <- err
				n.logger.Printf("Failed to dispatched id:%s\n", notifier.Id)
			} else {
				n.logger.Printf("Successfully dispatched id:%s\n", notifier.Id)
			}
		}(notifier)
	}
}

func (n *notificationsDispatcher) accumulateErrors(ch chan error) []error {
	errs := make([]error, 0)

	for {
		err := <-ch
		if err == nil {
			break
		}

		errs = append(errs, err)
	}

	return errs
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
	}
}
