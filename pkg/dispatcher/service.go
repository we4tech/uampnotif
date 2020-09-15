package dispatcher

import (
	"github.com/we4tech/uampnotif/pkg/integrations"
	"github.com/we4tech/uampnotif/pkg/notifications"
)

type dispatcher struct {
	notifier         *notifications.Notifier
	integrationSpecs map[string]integrations.Spec
}

//
// NewService returns an instance of dispatcher.
//
func NewService(notifier notifications.Notifier) *dispatcher {
	return &dispatcher{notifier: &notifier}
}

//
// LoadIntegrationSpecs iterates over the directory to load the integration
// plugins.
//
func (s *dispatcher) LoadIntegrationSpecs(configPath string) error {
	return nil
}

//
// Dispatch the notification based on the configuration and raises error
// when encounters.
//
func (s *dispatcher) Dispatch() error {
	return nil
}
