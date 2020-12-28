package receivers

import "testing"

func TestNewIntegrationSpec_ShouldReturnsIntegrationSpec(t *testing.T) {
	yamlConfig := `
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
`
	spec, err := NewSpec([]byte(yamlConfig))

	if err != nil {
		t.Errorf("could not parse the config. Error - %s", err)
	}

	if spec == nil {
		t.Error("could not find Spec")
	}

	if len(spec.Request.ValidHttpCodes) != 2 {
		t.Error()
	}

	if spec.Request.UrlTmpl == "" {
		t.Error()
	}

	if spec.Request.Method != "POST" {
		t.Error()
	}

	if spec.Request.BodyTmpl == "" {
		t.Error()
	}

	if spec.Request.Params.IsEmpty() {
		t.Error()
	}

	if spec.Request.Headers.IsEmpty() {
		t.Error()
	}

	if spec.Name != "Test" {
		t.Error()
	}

	if spec.Id != "test" {
		t.Error()
	}
}
