package main

import (
	"flag"
	"os"
)

type cliOpts struct {
	NotifierFile string
	ConfigDir    string
}

func (o *cliOpts) isEmpty() bool {
	return o.ConfigDir == "" || o.NotifierFile == ""
}

var parseFlags = func() *cliOpts {
	opts := &cliOpts{}

	flag.StringVar(&opts.NotifierFile, "n", "", "(Required) Locate notifiers.yml file")
	flag.StringVar(&opts.ConfigDir, "d", "", "(Required) Locate notifiers config directory")

	flag.Parse()

	if opts.isEmpty() {
		flag.Usage()

		os.Exit(1)
	}

	return opts
}
