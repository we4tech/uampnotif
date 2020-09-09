package controller

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

func Logging(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}
