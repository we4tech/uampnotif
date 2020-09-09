package templates

import (
	"bytes"
	"html/template"
	"time"

	"github.com/WeConnect/hello-tools/uampnotif/pkg/common_errors"
)

var functionMap = template.FuncMap{
	"now": time.Now,
}

//
// ExecuteTemplate a template and return the rendered respnse or an error.
//
func ExecuteTemplate(tmplName string, tmpl GoTmpl, ctx *TemplateContext) (string, error) {
	var parsedTmpl, err = template.
		New(tmplName).
		Funcs(functionMap).
		Parse(string(tmpl))

	if err != nil {
		return "", common_errors.TemplateParsingError{Err: err, Template: string(tmpl)}
	}

	buf := bytes.NewBuffer([]byte{})

	if err := parsedTmpl.Execute(buf, ctx); err != nil {
		return "", common_errors.TemplateParsingError{Err: err, Template: string(tmpl)}
	}

	return buf.String(), nil
}
