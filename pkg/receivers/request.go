package receivers

import (
	"github.com/we4tech/uampnotif/pkg/common_errors"
	"github.com/we4tech/uampnotif/pkg/templates"
)

type request struct {
	Params  *Params
	Headers *Headers
	Method  string

	ValidHttpCodes []int `yaml:"valid_http_codes"`

	UrlTmpl   templates.GoTmpl `yaml:"url_tmpl"`
	parsedUrl string

	BodyTmpl   templates.GoTmpl `yaml:"body_tmpl"`
	parsedBody string
}

func (r *request) Url(ctx *templates.TemplateContext) (string, error) {
	if r.parsedUrl != "" {
		return r.parsedUrl, nil
	}

	value, err := templates.ExecuteTemplate("urlTmpl", r.UrlTmpl, ctx)

	if err != nil {
		return "", common_errors.TemplateParsingError{
			Err:      err,
			Template: string(r.UrlTmpl),
		}
	}

	r.parsedUrl = value

	return r.parsedUrl, nil
}

func (r *request) Body(ctx *templates.TemplateContext) (string, error) {
	if r.parsedBody != "" {
		return r.parsedBody, nil
	}

	value, err := templates.ExecuteTemplate("bodyTmpl", r.BodyTmpl, ctx)

	if err != nil {
		return "", common_errors.TemplateParsingError{
			Err:      err,
			Template: string(r.UrlTmpl),
		}
	}

	r.parsedBody = value

	return r.parsedBody, nil
}
