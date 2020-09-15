package dispatcher

import (
	"github.com/we4tech/uampnotif/pkg/notifications"
	"os"
	"path"
	"testing"
)

func TestNewNotificationDispatcher(t *testing.T) {
	cfg := getNotificationConfig()

	dispatcher := NewNotificationDispatcher(cfg)

	if dispatcher == nil {
		t.Error()
	}
}

func TestDispatcher_Dispatch(t *testing.T) {
	// TODO: Mock http call and create a stubbed notification config
}

func getNotificationConfig() *notifications.Config {
	dir, _ := os.Getwd()
	configFile := path.Join(dir, "../../config/notifications.yml")

	cfgParser := notifications.NewDefaultConfigParser()
	if cfg, err := cfgParser.Read(configFile); err != nil {
		panic(err)
	} else {
		return cfg
	}
	return nil
}
