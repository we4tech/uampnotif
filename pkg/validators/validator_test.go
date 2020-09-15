package validators

import (
	"github.com/we4tech/uampnotif/pkg/templates"
	"testing"

	"github.com/we4tech/uampnotif/pkg/integrations"
)

var IntegrationTestTables = []struct {
	testName         string
	integration      *integrations.Spec
	expectValid      bool
	expectedErrorsOn []string
}{
	{"should be an invalid integration without name",
		buildIntegration(""),
		false,
		[]string{"name"}},
	{"should be an invalid integration without valid_http_codes",
		buildIntegration(""),
		false,
		[]string{"request.validHttpCodes"}},
	{"should be an invalid integration without id",
		buildIntegration(""),
		false,
		[]string{"id"}},
	{"should be an invalid integration without request.method",
		buildIntegration(""),
		false,
		[]string{"request.method"}},
	{"should be an invalid integration without request.url",
		buildIntegration(""),
		false,
		[]string{"request.urlTmpl"}},
	{"should be a valid integration",
		buildIntegration(`
name: Test
id: test
request:
  valid_http_codes:
    - 201
    - 200
  params:
    - name: app_id
      label: App ID
    - name: api_key
      label: API Key
  method: POST
  url_tmpl: 'https://api.test.local/v2/applications/{{.FindParam "app_id"}}/deployments.json'
  headers:
    - name: Content-Type
      value_tmpl: application/json
    - name: X-Api-Key
      value_tmpl: '{{.FindParam "api_key"}}'
  body_tmpl: |
    { "deployment": { "revision": "{{.FindEnv "commit_hash"}}" } }
`),
		true,
		[]string{}},
}

func buildIntegration(config string) *integrations.Spec {
	spec, _ := integrations.NewSpec([]byte(config))

	return spec
}

func TestValidateIntegration(t *testing.T) {
	for _, tt := range IntegrationTestTables {
		t.Run(tt.testName, func(t *testing.T) {
			params := map[string]string{}
			envVars := map[string]string{}
			tmplCtx := templates.NewTemplateContext(params, envVars)
			validator := NewValidator(tt.integration, tmplCtx)
			valid := validator.Validate()
			errors := validator.GetErrors()

			if tt.expectValid {
				if !valid {
					t.Errorf("expected to be valid but received errors on %s", errors)
				}
			} else {
				if valid {
					t.Error("expected not to be valid")
				}

				for _, field := range tt.expectedErrorsOn {
					if !errors.HasError(field) {
						t.Errorf("could not find error on %s field", field)
					} else if errors.GetError(field) != "is required" {
						t.Errorf("could not find error message 'is required'")
					}
				}
			}
		})
	}
}
