package notifiers

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/we4tech/uampnotif/pkg/common_errors"
)

func TestReadShouldRaiseConfigNotFound(t *testing.T) {
	parser := NewParser()
	_, err := parser.Read("../invalid-file-path.yaml")

	if err == nil {
		t.Error("should raise error")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(common_errors.ConfigNotFound{}) {
		t.Error("should raise ConfigNotFoundError")
	}
}

func TestReadInternalShouldRaiseConfigParsingError(t *testing.T) {
	invalidYaml := []byte("settings: hello: hello")
	invalidYamlFile := "a-file-name-does-matter.yml"

	parser := &parser{}

	_, err := parser.readInternal(invalidYaml, invalidYamlFile)

	if err == nil {
		t.Error("should raise error")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(common_errors.ConfigParsingError{}) {
		t.Error("should raise ConfigParsingError")
	}
}

func TestReadShouldParseConfig(t *testing.T) {
	dir, _ := os.Getwd()
	configFile := path.Join(dir, "../../config/notifiers.yml")

	parser := NewParser()
	notifiers, err := parser.Read(configFile)

	if err != nil {
		t.Errorf("could not parse notifiers.yml. Error - %s", err)
	}

	validateNotifiers(notifiers, t)
}

func validateNotifiers(notifiers *Config, t *testing.T) {
	t.Run("should have default settings", func(t *testing.T) {
		if notifiers.DefaultSettings.Retries != 3 {
			t.Error("could not find retries == 3")
		}

		if !notifiers.DefaultSettings.Async {
			t.Error("could not find async == true")
		}

		if notifiers.DefaultSettings.OnError != "ignore" {
			t.Error("could not find on_error: ignore")
		}

		if len(notifiers.DefaultSettings.OnErrorNotifiers) != 1 {
			t.Error("could not find len(on_error_notifiers) == 1")
		}
	})

	expectedNotifiers := []string{"newrelic", "rollbar", "slack", "sox-auditor"}

	for _, notifId := range expectedNotifiers {
		t.Run(
			fmt.Sprintf("should find notifier - %s", notifId),
			func(t *testing.T) {
				_, found := findNotifier(notifId, notifiers)

				if !found {
					t.Errorf("could not find - %s", notifId)
				}
			})

		t.Run("should have parameters", func(t *testing.T) {
			n, _ := findNotifier(notifId, notifiers)

			if n.Params.IsEmpty() {
				t.Errorf("could not find parameters")
			}
		})
	}

	t.Run("should have settings for sox-auditor", func(t *testing.T) {
		n, _ := findNotifier("sox-auditor", notifiers)

		if n.Settings.OnError != "fatal" {
			t.Errorf("could not find setting.on_error")
		}
	})
}

func findNotifier(name string, config *Config) (Notifier, bool) {
	found := false
	var notifier Notifier

	for _, n := range config.Notifiers {
		if n.Id == name {
			found = true
			notifier = n
		}
	}

	return notifier, found
}
