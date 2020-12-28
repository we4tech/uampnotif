package notifcfg

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
	configFile := path.Join(dir, "../../test-configs/notification.yml")

	parser := NewParser()
	cfg, err := parser.Read(configFile)

	if err != nil {
		t.Errorf("could not parse notification.yml. Error - %s", err)
	}

	validateConfig(cfg, t)
}

func validateConfig(cfg *Config, t *testing.T) {
	t.Run("should have default settings", func(t *testing.T) {
		if cfg.DefaultSettings.Retries != 3 {
			t.Error("could not find retries == 3")
		}

		if !cfg.DefaultSettings.Async {
			t.Error("could not find async == true")
		}

		if cfg.DefaultSettings.OnError != "ignore" {
			t.Error("could not find on_error: ignore")
		}

		if len(cfg.DefaultSettings.OnErrorReceivers) != 1 {
			t.Error("could not find len(on_error_receivers) == 1")
		}
	})

	expectedReceivers := []string{"newrelic", "rollbar", "slack", "sox-auditor"}

	for _, receiverId := range expectedReceivers {
		t.Run(
			fmt.Sprintf("should find receiver - %s", receiverId),
			func(t *testing.T) {
				_, found := findReceiver(receiverId, cfg)

				if !found {
					t.Errorf("could not find - %s", receiverId)
				}
			})

		t.Run("should have parameters", func(t *testing.T) {
			n, _ := findReceiver(receiverId, cfg)

			if n.Params.IsEmpty() {
				t.Errorf("could not find parameters")
			}
		})
	}

	t.Run("should have settings for sox-auditor", func(t *testing.T) {
		n, _ := findReceiver("sox-auditor", cfg)

		if n.Settings.OnError != "fatal" {
			t.Errorf("could not find setting.on_error")
		}
	})
}

func findReceiver(name string, config *Config) (Receiver, bool) {
	found := false
	var receiver Receiver

	for _, n := range config.Receivers {
		if n.Id == name {
			found = true
			receiver = n
		}
	}

	return receiver, found
}
