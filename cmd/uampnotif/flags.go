package main

import (
	"flag"
	"os"
)

type cliOpts struct {
	NotificationCfgFile string
	ReceiverSpecDir     string
}

func (o *cliOpts) isEmpty() bool {
	return o.ReceiverSpecDir == "" || o.NotificationCfgFile == ""
}

var parseFlags = func() *cliOpts {
	opts := &cliOpts{}

	flag.StringVar(&opts.NotificationCfgFile, "n", "", "(Required) Locate notification2.yml file")
	flag.StringVar(&opts.ReceiverSpecDir, "d", "", "(Required) Locate receiver (*.yml) specs directory")

	flag.Parse()

	if opts.isEmpty() {
		flag.Usage()

		os.Exit(1)
	}

	return opts
}
