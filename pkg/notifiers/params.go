package notifiers

import (
	"github.com/WeConnect/hello-tools/uampnotif/pkg/common_errors"
	"github.com/WeConnect/hello-tools/uampnotif/pkg/templates"
)

//
// Params contains a map of key-value pair.
//
type Params map[string]string

//
// IsEmpty returns true if len(params) == 0.
//
func (p *Params) IsEmpty() bool {
	return len(*p) == 0
}

func (p *Params) GetValue(ctx *templates.TemplateContext, key string) (string, error) {
	value, ok := (*p)[key]

	if !ok {
		return "", common_errors.KeyNotFoundError{Key: key}
	}

	value, err := templates.ExecuteTemplate(
		"Params.Value", templates.GoTmpl(value), ctx)

	if err != nil {
		return "", err
	}

	return value, nil
}
