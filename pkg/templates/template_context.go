package templates

import (
	"github.com/WeConnect/hello-tools/uampnotif/pkg/common_errors"
)

//
// TemplateContext stores all contextual variables injected from the system and
// integration data store. ie - callback_url (data store), commit_hash (system var)
//
type TemplateContext struct {
	Params map[string]string
	Env    map[string]string
}

//
// NewTemplateContext returns an instance of TemplateContext.
//
func NewTemplateContext(
	params map[string]string,
	envVars map[string]string) *TemplateContext {
	return &TemplateContext{Params: params, Env: envVars}
}

//
// FindParam returns a parameter value if exists otherwise returns an error.
//
func (ctx *TemplateContext) FindParam(needle string) (string, error) {
	if value, ok := ctx.Params[needle]; ok {
		return value, nil
	}

	return "", common_errors.KeyNotFoundError{Key: needle}
}

//
// ParamExists returns true if found a matching parameter.
//
func (ctx *TemplateContext) ParamExists(needle string) bool {
	_, ok := ctx.Params[needle]

	return ok
}

//
// FindEnv returns the environment value by the specified key.
//
func (ctx *TemplateContext) FindEnv(key string) (string, error) {
	if value, ok := ctx.Env[key]; ok {
		return value, nil
	}

	return "", common_errors.KeyNotFoundError{Key: key}
}

//
// EnvExists returns true if the environment key exists.
//
func (ctx *TemplateContext) EnvExists(key string) bool {
	_, ok := ctx.Env[key]

	return ok
}
