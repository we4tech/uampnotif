package dispatcher

import (
	"context"
	"fmt"
	"github.com/we4tech/uampnotif/pkg/clients"
	"github.com/we4tech/uampnotif/pkg/notifcfg"
	"github.com/we4tech/uampnotif/pkg/receivers"
	"golang.org/x/sync/errgroup"
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
	notificationCfg *notifcfg.Config
	specCfg         map[string]*receivers.Spec
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
// Dispatch triggers events for all notifcfg using separate go-routine.
//
// TODO(HK): Add support for SYNC requests.
func (n *notificationsDispatcher) Dispatch(ctx context.Context) error {
	defer n.closeChannels()

	if err := n.dispatchInAsync(ctx); err != nil {
		n.logger.Println("Failed to dispatch all notifications")

		return &dispatchError{Errors: []error{err}}
	}

	fmt.Println("Dispatched")
	fmt.Println("Closed err")

	n.logger.Println("Successfully dispatched all notifications")

	return nil
}

func (n *notificationsDispatcher) SetMockClient(client clients.ClientImpl) {
	n.mockClient = client
}

// TODO(HK): Use context
func (n *notificationsDispatcher) dispatchNotification(
	_ context.Context,
	receiver notifcfg.Receiver,
) error {
	var (
		retries    = 0
		response   *clients.Response
		maxRetries = n.maxRetries(receiver.Settings)
	)

	n.trigger(InTransit, receiver, retries, nil, nil)

restart:
	client, err := clients.NewHttpRequest(n.specCfg[receiver.Id], receiver.Params, n.envVars, n.logger)
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
		n.trigger(Success, receiver, retries, response, nil)
		return nil
	}

	if retries < maxRetries {
		n.trigger(Retrying, receiver, retries, response, nil)

		retries++

		goto restart
	}

errorHandler:
	n.trigger(Error, receiver, retries, response, err)

	return err
}

func (n *notificationsDispatcher) trigger(
	state DispatchState,
	receiver notifcfg.Receiver,
	retries int,
	response *clients.Response,
	err error,
) {
	if n.events == nil {
		return
	}

	n.events <- DispatchEvent{
		State:        state,
		ReceiverId:   receiver.Id,
		ReceiverDesc: receiver.Desc,
		Retries:      retries,
		Response:     response,
		Error:        err,
	}
}

func (n *notificationsDispatcher) maxRetries(settings *notifcfg.Setting) int {
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

func (n *notificationsDispatcher) dispatchInAsync(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	for _, receiver := range n.notificationCfg.Receivers {
		n.logger.Printf("Dispatching [%s:%s]\n", receiver.Id, receiver.Desc)

		g.Go(func() error {
			if err := n.dispatchNotification(ctx, receiver); err != nil {
				n.logger.Printf("Failed to dispatched [%s:%s]\n", receiver.Id, receiver.Desc)

				return err
			}

			n.logger.Printf("Successfully dispatched [%s:%s]\n", receiver.Id, receiver.Desc)

			return nil
		})
	}

	return g.Wait()
}

func (n *notificationsDispatcher) closeChannels() {
	if n.events != nil {
		fmt.Println("Closing events")
		close(n.events)
	}

	if n.done != nil {
		close(n.done)
	}
}

//
// NewNotificationDispatcher returns an implementation of *Dispatcher* service.
//
func NewNotificationDispatcher(
	logger *log.Logger,
	specs map[string]*receivers.Spec,
	config *notifcfg.Config,
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
