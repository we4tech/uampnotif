package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var commitHash = "adc83b19e793491b1c6ea0fd8b46cd9f32e592fc"

func TestMainWithHttpServer(t *testing.T) {
	opts := &cliOpts{
		"../../config/testconfigs/notifiers.yml",
		"../../config/testconfigs/configs",
	}

	server := runServer(t)
	defer server.Close()

	origFlags := parseFlags
	defer func() { parseFlags = origFlags }()

	parseFlags = func() *cliOpts { return opts }

	_ = os.Setenv("NOTIF_SERVER_URL", server.URL)
	_ = os.Setenv("COMMIT_HASH", commitHash)

	main()
}

func runServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		log.Printf("Received request with payload: %s\n", body)

		if !strings.Contains(string(body), commitHash) {
			t.Error("response body doesn't contain the commit hash env-var.")
		}

		_, err := fmt.Fprintf(w, "hello world")
		if err != nil {
			panic(err)
		}
	}))
}
