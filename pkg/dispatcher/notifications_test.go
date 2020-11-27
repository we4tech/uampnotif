package dispatcher_test

import (
	"context"
	"github.com/we4tech/uampnotif/pkg/dispatcher"
	"github.com/we4tech/uampnotif/pkg/integrations"
	"github.com/we4tech/uampnotif/pkg/notifications"
	"github.com/we4tech/uampnotif/pkg/testutils"
	"log"
	"os"
	"path"
	"testing"
)

func TestDispatcher_Dispatch(t *testing.T) {
	cfg := getNotificationConfig()
	ctx := context.Background()
	specsMap := map[string]*integrations.Spec{
		"test": getIntegrationSpec(),
	}
	envVars := map[string]string{
		"commit_hash": "e983f374794de9c64e3d1c1de1d490c0756eeeff",
	}
	params := map[string]string{
		"app_id":  "app1234",
		"api_key": "abcdef",
	}

	t.Run("dispatches without any error", func(t *testing.T) {
		mockClient := testutils.NewMockHttpClient(200, []byte("[]"))

		d := dispatcher.NewNotificationDispatcher(specsMap, cfg, params, envVars)
		d.SetMockClient(mockClient)

		if err := d.Dispatch(ctx); err != nil {
			t.Error("[Test] Error: ", err)
		}
	})

	t.Run("dispatches events", func(t *testing.T) {
		mockClient := testutils.NewMockHttpClient(200, []byte("[]"))

		d := dispatcher.NewNotificationDispatcher(specsMap, cfg, params, envVars)
		d.SetMockClient(mockClient)

		go func() { _ = d.Dispatch(ctx) }()

		var lastEvent dispatcher.DispatchEvent
		inTransitCount, successCount, errorCount := 0, 0, 0

		for {
			lastEvent = <-d.Channel()

			log.Printf("Rcv: %+v", lastEvent)

			if lastEvent.State == dispatcher.InTransit {
				inTransitCount++
			} else if lastEvent.State == dispatcher.Error {
				t.Error("expected success")
				errorCount++
			} else if lastEvent.State == dispatcher.Success {
				log.Println("Successfully dispatched")
				successCount++
			}

			if successCount+errorCount > 1 {
				break
			}
		}

		if inTransitCount != 2 {
			t.Error()
		}

		if errorCount != 0 {
			t.Error()
		}

		if successCount != 2 {
			t.Error()
		}
	})
}

func getIntegrationSpec() *integrations.Spec {
	dir, _ := os.Getwd()
	configFile := path.Join(dir, "../../config/integrations/test.yml")

	if spec, err := integrations.NewConfigParser().Read(configFile); err != nil {
		panic(err)
	} else {
		return spec
	}
}

func getNotificationConfig() *notifications.Config {
	dir, _ := os.Getwd()
	configFile := path.Join(dir, "../../config/test-fixtures/notifs-simple.yml")

	cfgParser := notifications.NewConfigParser()
	if cfg, err := cfgParser.Read(configFile); err != nil {
		panic(err)
	} else {
		return cfg
	}
	return nil
}
