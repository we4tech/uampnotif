package main

import (
	"context"
	"github.com/we4tech/uampnotif/pkg/configs"
	"github.com/we4tech/uampnotif/pkg/dispatcher"
	"github.com/we4tech/uampnotif/pkg/notifiers"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "[uampnotif] ", log.LUTC&log.Lmsgprefix&log.Lshortfile)
}

func main() {
	opts := parseFlags()

	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	logger.Println("Loading configuration...")

	parser := notifiers.NewParser()
	notifierFile := path.Join(wd, opts.NotifierFile)

	config, err := parser.Read(notifierFile)
	if err != nil {
		log.Panicf("could not parse file from %s. error: %s", notifierFile, err)
	}

	specsMap := buildSpecsMap(opts, wd)
	params := make(map[string]string)
	envVars := buildEnvVarsMap()
	pCtx := context.Background()
	ctx, cancel := context.WithCancel(pCtx)

	logger.Println("Preparing dispatcher")

	d := dispatcher.NewNotificationDispatcher(logger, specsMap, config, params, envVars)

	go monitorEvents(ctx, logger, d)

	if err := d.Dispatch(ctx); err != nil {
		log.Panicf("failed to dispatch successfully. error: %s", err)
	}

	cancel()
}

func monitorEvents(ctx context.Context, logger *log.Logger, d dispatcher.Dispatcher) {
	for {
		select {
		case <-ctx.Done():
			goto returnCtl
		case event := <-d.Events():
			logger.Printf("Event: %+v\n", event)
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

func buildSpecsMap(opts *cliOpts, wd string) map[string]*configs.Spec {
	specsMap := make(map[string]*configs.Spec)

	fullPath := path.Join(wd, opts.ConfigDir)

	configFiles, err := ioutil.ReadDir(fullPath)
	if err != nil {
		log.Panicf("could not read directory. error: %s", err)
	}

	parser := configs.NewParser()

	for _, file := range configFiles {
		fullPath := path.Join(wd, opts.ConfigDir, file.Name())
		spec, err := parser.Read(fullPath)
		if err != nil {
			log.Panicf("could not load config file - %s. error: %s", fullPath, err)
		}

		specsMap[spec.Id] = spec
	}

	return specsMap
}
