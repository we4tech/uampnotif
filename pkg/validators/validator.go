package validators

import (
	"github.com/WeConnect/hello-tools/uampnotif/pkg/integrations"
	"github.com/WeConnect/hello-tools/uampnotif/pkg/templates"
)

//
// ## WHY?:
//
// The integration-validator validates a parsed struct to ensure all required
// attributes exist.
//

//
// integrationSpecValidator implements a Validator to validate an IntegrationSpec struct.
//
type integrationSpecValidator struct {
	templateContext *templates.TemplateContext
	integration     *integrations.IntegrationSpec
	errors          *validationErrors
	valid           bool
}

//
// validationErrors keeps all validation specific errors.
//
type validationErrors map[string]string

//
// HasError returns true if a validation error exists
//
func (ve *validationErrors) HasError(field string) bool {
	return (*ve)[field] != ""
}

//
// GetError returns the error string of the specified field.
//
func (ve *validationErrors) GetError(field string) string {
	return (*ve)[field]
}

//
// GetErrors returns the list of validation errors.
//
func (iv *integrationSpecValidator) GetErrors() ValidationErrors {
	return iv.errors
}

//
// Validate returns true if the integration configuration is valid.
//
func (iv *integrationSpecValidator) Validate() bool {
	iv.errors = &validationErrors{}

	if iv.integration.Id == "" {
		iv.addError("id", "is required")
	}

	if iv.integration.Name == "" {
		iv.addError("name", "is required")
	}

	if len(iv.integration.Request.ValidHttpCodes) == 0 {
		iv.addError("request.validHttpCodes", "is required")
	}

	if iv.integration.Request.UrlTmpl == "" {
		iv.addError("request.urlTmpl", "is required")
	}

	if iv.integration.Request.Method == "" {
		iv.addError("request.method", "is required")
	}

	if len(*iv.errors) == 0 {
		iv.valid = true
	}

	return iv.valid
}

func (iv *integrationSpecValidator) addError(field string, message string) {
	(*iv.errors)[field] = message
}

//
// NewValidator returns an instance of validator.
//
func NewValidator(
	integration *integrations.IntegrationSpec,
	templateCtx *templates.TemplateContext) Validator {
	return &integrationSpecValidator{
		integration:     integration,
		templateContext: templateCtx}
}
