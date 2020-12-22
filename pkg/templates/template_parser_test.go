package templates

import (
	"regexp"
	"testing"
)

func TestExecuteTemplateShouldReturnString(t *testing.T) {
	params := map[string]string{
		"host": "example.org",
		"path": "hello-world",
	}
	ctx := TemplateContext{
		Params: params,
	}
	tpl := "http://{{.FindParam \"host\"}}/{{.FindParam \"path\"}}"

	value, err := ExecuteTemplate(
		"exec-test-tmpl", GoTmpl(tpl), &ctx)

	if err != nil {
		t.Errorf("could not parse the template. Error - %s", err)
	}

	if value != "http://example.org/hello-world" {
		t.Errorf("could not match expected url")
	}
}

func TestExecuteTemplateShouldRaiseErrorWithInvalidTemplate(t *testing.T) {
	params := map[string]string{
		"host": "example.org",
		"path": "hello-world",
	}
	ctx := TemplateContext{
		Params: params,
	}
	tmpl := "Hi {{Hello}}"

	_, err := ExecuteTemplate(
		"exec-failure-test", GoTmpl(tmpl), &ctx)

	if err == nil {
		t.Error("could not find template parsing error")
	}

	methodNotDefined := regexp.MustCompile("function \"Hello\" not defined")

	if !methodNotDefined.MatchString(err.Error()) {
		t.Errorf("could not find expected error message - %s", err)
	}
}
