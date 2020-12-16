package dispatcher

import (
	"github.com/we4tech/uampnotif/pkg/configs"
	"github.com/we4tech/uampnotif/pkg/notifiers"
)

type dispatcher struct {
	notifier         *notifiers.Notifier
	integrationSpecs map[string]configs.Spec
}

//
// NewService returns an instance of dispatcher.
//
func NewService(notifier notifiers.Notifier) *dispatcher {
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
