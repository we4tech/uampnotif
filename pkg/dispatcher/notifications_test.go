package dispatcher_test

import (
	"context"
	"fmt"
	"github.com/we4tech/uampnotif/pkg/dispatcher"
	"github.com/we4tech/uampnotif/pkg/notifcfg"
	"github.com/we4tech/uampnotif/pkg/receivers"
	"github.com/we4tech/uampnotif/pkg/testutils"
	"log"
	"os"
	"path"
	"testing"
)

var logger = log.New(os.Stdout, "[testbin] ", log.Lshortfile)

func TestDispatcher_Dispatch(t *testing.T) {
	cfg := getNotificationConfig()
	ctx := context.Background()
	specsMap := map[string]*receivers.Spec{
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

		d := dispatcher.NewNotificationDispatcher(logger, specsMap, cfg, params, envVars)
		d.SetMockClient(mockClient)

		go func() {
			for {
				event := <-d.Events()
				if event.ReceiverId == "" {
					break
				}

				fmt.Printf("Event: %+v\n", event)
			}
		}()

		if err := d.Dispatch(ctx); err != nil {
			t.Error("[Test] Error: ", err)
		}
	})

	t.Run("dispatches events", func(t *testing.T) {
		ctx := context.Background()
		mockClient := testutils.NewMockHttpClient(200, []byte("[]"))

		d := dispatcher.NewNotificationDispatcher(logger, specsMap, cfg, params, envVars)
		d.SetMockClient(mockClient)

		go func() { _ = d.Dispatch(ctx) }()

		var lastEvent dispatcher.DispatchEvent
		inTransitCount, successCount, errorCount := 0, 0, 0

		for {
			select {
			case <-d.Done():
				fmt.Println("Done test")
				goto continueTest

			case lastEvent = <-d.Events():
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
			}
		}

	continueTest:
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

func getIntegrationSpec() *receivers.Spec {
	dir, _ := os.Getwd()
	configFile := path.Join(dir, "../../test-configs/receivers/test.yml")

	if spec, err := receivers.NewParser().Read(configFile); err != nil {
		panic(err)
	} else {
		return spec
	}
}

func getNotificationConfig() *notifcfg.Config {
	dir, _ := os.Getwd()
	configFile := path.Join(dir, "../../test-configs/notification2.yml")

	cfgParser := notifcfg.NewParser()

	cfg, err := cfgParser.Read(configFile)
	if err != nil {
		panic(err)
	}

	return cfg
}
