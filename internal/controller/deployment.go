package controller

import (
	"net/http"
)

type Deployment struct {
}

const HttpOk = 200

func NewDeployment() Deployment {
	return Deployment{}
}

func (d Deployment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(HttpOk)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("Hello world"))
}
