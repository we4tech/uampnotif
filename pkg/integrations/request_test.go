package integrations

import (
	"github.com/we4tech/uampnotif/pkg/templates"
	"testing"
)

func TestUrlShouldConstructWithParam(t *testing.T) {
	expectedUrl := "https://example.org/api/endpoint"

	req := request{UrlTmpl: "{{.FindParam \"schema\"}}://" +
		"{{.FindParam \"host\"}}/{{.FindParam \"path\"}}"}
	ctx := templates.TemplateContext{Params: map[string]string{
		"schema": "https",
		"host":   "example.org",
		"path":   "api/endpoint",
	}}

	if url, _ := req.Url(&ctx); url != expectedUrl {
		t.Errorf("could not find url as expected - %s", expectedUrl)
	}
}
