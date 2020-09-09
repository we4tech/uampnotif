package integrations

import (
	"github.com/WeConnect/hello-tools/uampnotif/pkg/templates"
	"log"
)

//
// Headers is an alias for an array of header
type Headers []header

//
// ForEach implements Iterator method to provide access to internal data
// structure.
//
func (h *Headers) ForEach(cb func(i int, h ParsedHeader) (bool, error)) error {
	var loopError error

	if h.IsEmpty() {
		log.Printf("IterableHeaders.ForEach: Empty Headers set")

		return nil
	}

	for i, header := range *h {
		if loopBreak, err := cb(i, &header); loopBreak {
			loopError = err

			break
		}
	}

	return loopError
}

func (h *Headers) IsEmpty() bool {
	return h == nil || len(*h) == 0
}

type header struct {
	Name        string
	ValueTmpl   templates.GoTmpl `yaml:"value_tmpl"`
	parsedValue string
}

func (h *header) GetValue(ctx *templates.TemplateContext) (string, error) {
	value, err := templates.ExecuteTemplate("header.valueTmpl", h.ValueTmpl, ctx)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (h *header) GetName() string {
	return h.Name
}
