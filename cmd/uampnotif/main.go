package main

import (
	"github.com/WeConnect/hello-tools/uampnotif/internal/controller"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"

	"log"
	"os"
)

//
// Global variables.
//

var (
	// HTTP binding address. By default binds to all interfaces.
	address = "0.0.0.0:3030"
)

func init() {
	if addrs, found := os.LookupEnv("BIND_ADDRS"); found {
		log.Printf("Setting BIND_ADDRS from env var - %s", addrs)

		address = addrs
	}
}

func main() {
	r := mux.NewRouter()
	r.Handle(
		"/deployments",
		handlers.LoggingHandler(
			os.Stdout, controller.NewDeployment()))

	log.Printf("Server started at %s", address)

	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Failed to launch with error %s", err)
	}
}
