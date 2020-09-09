package dispatcher

import (
	"github.com/WeConnect/hello-tools/uampnotif/pkg/integrations"
	"github.com/WeConnect/hello-tools/uampnotif/pkg/notifiers"
)

type service struct {
	notifier *notifiers.Notifier
	integrationSpecs map[string]integrations.IntegrationSpec
}

//
// NewService returns an instance of service.
//
func NewService(notifier notifiers.Notifier) *service {
	return &service{notifier: &notifier}
}

//
// LoadIntegrationSpecs iterates over the directory to load the integration
// plugins.
//
func (s *service) LoadIntegrationSpecs(configPath string) error {
	return nil
}

//
// Dispatch the notification based on the configuration and raises error
// when encounters.
//
func (s *service) Dispatch() error {
	return nil
}
