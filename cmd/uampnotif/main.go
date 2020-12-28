package main

import (
	"context"
	"github.com/gookit/color"
	"github.com/we4tech/uampnotif/pkg/dispatcher"
	"github.com/we4tech/uampnotif/pkg/notifcfg"
	"github.com/we4tech/uampnotif/pkg/receivers"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var logger *log.Logger

func init() {
	logger = log.New(
		os.Stdout,
		color.Gray.Sprint("[uampnotif] "),
		log.Ldate|log.Lshortfile,
	)
}

func main() {
	opts := parseFlags()

	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	logger.Println("Loading configuration...")

	parser := notifcfg.NewParser()
	cfgFile := path.Join(wd, opts.NotificationCfgFile)

	config, err := parser.Read(cfgFile)
	if err != nil {
		log.Panicf("could not parse file from %s. error: %s", cfgFile, err)
	}

	specsMap := buildSpecsMap(opts, wd)
	params := make(map[string]string)
	envVars := buildEnvVarsMap()
	pCtx := context.Background()
	ctx, cancel := context.WithCancel(pCtx)

	d := dispatcher.NewNotificationDispatcher(logger, specsMap, config, params, envVars)

	go monitorEvents(ctx, logger, d)

	if err := d.Dispatch(ctx); err != nil {
		logger.Panicf(color.Red.Sprintf("Failed to dispatch successfully. error: %s", err))
	}

	cancel()
}

func monitorEvents(ctx context.Context, logger *log.Logger, d dispatcher.Dispatcher) {
	for {
		select {
		case <-ctx.Done():
			goto returnCtl
		case event, ok := <-d.Events():
			if !ok {
				goto returnCtl
			}
			logger.Print(color.Yellow.Sprintf("Event: %+v\n", event))
		}
	}

returnCtl:
	return
}

func buildEnvVarsMap() map[string]string {
	envVars := make(map[string]string)

	for _, envLine := range os.Environ() {
		parts := strings.Split(envLine, "=")

		envVars[parts[0]] = parts[1]
	}

	return envVars
}

func buildSpecsMap(opts *cliOpts, wd string) map[string]*receivers.Spec {
	specsMap := make(map[string]*receivers.Spec)

	fullPath := path.Join(wd, opts.ReceiverSpecDir)

	configFiles, err := ioutil.ReadDir(fullPath)
	if err != nil {
		log.Panicf("could not read directory. error: %s", err)
	}

	parser := receivers.NewParser()

	for _, file := range configFiles {
		fullPath := path.Join(wd, opts.ReceiverSpecDir, file.Name())
		spec, err := parser.Read(fullPath)
		if err != nil {
			log.Panicf("could not load config file - %s. error: %s", fullPath, err)
		}

		specsMap[spec.Id] = spec
	}

	return specsMap
}
