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
	"sync"
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

	logger.Println("Loading configuration...")

	parser := notifcfg.NewParser()
	cfgFile := opts.NotificationCfgFile

	config, err := parser.Read(cfgFile)
	if err != nil {
		log.Panicf("could not parse file from %s. error: %s", cfgFile, err)
	}

	specsMap := buildSpecsMap(opts)
	params := make(map[string]string)
	envVars := buildEnvVarsMap()
	pCtx := context.Background()
	wg := &sync.WaitGroup{}

	d := dispatcher.NewNotificationDispatcher(logger, specsMap, config, params, envVars)

	wg.Add(1)
	go monitorEvents(logger, d.Events(), wg)

	if err := d.Dispatch(pCtx); err != nil {
		logger.Println(color.Red.Sprint("Failed to dispatch successfully"))
	}

	wg.Wait()
}

func monitorEvents(logger *log.Logger, d chan dispatcher.DispatchEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		event, ok := <-d
		if !ok {
			break
		}
		logger.Print(color.Yellow.Sprintf("Event: %+v\n", event))
	}
}

func buildEnvVarsMap() map[string]string {
	envVars := make(map[string]string)

	for _, envLine := range os.Environ() {
		parts := strings.Split(envLine, "=")

		envVars[parts[0]] = parts[1]
	}

	return envVars
}

func buildSpecsMap(opts *cliOpts) map[string]*receivers.Spec {
	specsMap := make(map[string]*receivers.Spec)

	configFiles, err := ioutil.ReadDir(opts.ReceiverSpecDir)
	if err != nil {
		log.Panicf("could not read directory. error: %s", err)
	}

	parser := receivers.NewParser()

	for _, file := range configFiles {
		fullPath := path.Join(opts.ReceiverSpecDir, file.Name())
		spec, err := parser.Read(fullPath)
		if err != nil {
			log.Panicf("could not load config file - %s. error: %s", fullPath, err)
		}

		specsMap[spec.Id] = spec
	}

	return specsMap
}
