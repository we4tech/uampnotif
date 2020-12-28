package receivers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/we4tech/uampnotif/pkg/common_errors"
	"github.com/we4tech/uampnotif/pkg/templates"
)

func TestReadShouldRaiseConfigNotFound(t *testing.T) {
	_, err := NewParser().Read("invalid-path.yaml")

	if err == nil {
		t.Error("could not find ConfigNotFound error")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(common_errors.ConfigNotFound{}) {
		t.Error("could not find ConfigNotFound error")
	}
}

func TestReadInternalRaiseConfigParsingError(t *testing.T) {
	parser := NewParser()
	_, err := parser.ReadBytes([]byte("invalid_yaml: _invalid: _invalid"))

	if err == nil {
		t.Error("could not find error")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(common_errors.ConfigParsingError{}) {
		t.Error("could not find ConfigParsingError")
	}
}

func TestReadShouldReadConfigs(t *testing.T) {
	dir, _ := os.Getwd()
	rootPath := path.Join(dir, "../../test-configs/receivers")
	configFiles, err := ioutil.ReadDir(rootPath)

	if err != nil {
		t.Errorf(
			"could not find list of config files. Error - %s", err)
	}

	for _, file := range configFiles {
		t.Run(
			fmt.Sprintf("should parse %s", file.Name()),
			func(t *testing.T) {
				integration, err := getIntegration(path.Join(rootPath, file.Name()))

				if err != nil {
					t.Errorf("could not parse newrelic.yml - error - %s", err)
				}

				validateIntegration(integration, t)
			})
	}
}

func validateIntegration(integration *Spec, t *testing.T) {
	ctx := buildContext(integration)

	t.Run("should find name", func(t *testing.T) {
		if integration.Name == "" {
			t.Error("could not find name")
		}
	})

	t.Run("should find id", func(t *testing.T) {
		if integration.Id == "" {
			t.Error("could not find id")
		}
	})

	t.Run("should find request Params", func(t *testing.T) {
		if integration.Request.Params.IsEmpty() {
			t.Error("could not find request Params")
		}
	})

	t.Run("should find request Headers", func(t *testing.T) {
		if integration.Request.Headers.IsEmpty() {
			t.Error("could not find request Headers")
		}
	})

	t.Run("should find request method", func(t *testing.T) {
		if integration.Request.Method == "" {
			t.Error("could not find request method")
		}
	})

	t.Run("should find request body tmpl", func(t *testing.T) {
		if body, err := integration.Request.Body(ctx); err != nil {
			t.Errorf("could not parse bodyTmpl - error: %s", err)
		} else if body == "" {
			t.Error("could not find request body tmpl")
		}
	})

	t.Run("should find request url tmpl", func(t *testing.T) {
		if url, err := integration.Request.Url(ctx); err != nil {
			t.Errorf("could not find request url tmpl. Error - %s", err)
		} else if url == "" {
			t.Errorf("could not find value for request.Url. Error - %s", err)
		}
	})

	t.Run("should find request.Params", func(t *testing.T) {
		integration.Request.Params.ForEach(func(i int, p *Param) {
			if p.Name == "" {
				t.Error("could not find param.name")
			}

			if p.Label == "" {
				t.Error("could not find param.label")
			}
		})
	})

	t.Run("should find request.validHttpCodes", func(t *testing.T) {
		if len(integration.Request.ValidHttpCodes) == 0 {
			t.Error("could not find request.ValidHTTPCodes")
		}
	})

	t.Run("should find request.Headers", func(t *testing.T) {
		err := integration.Request.Headers.ForEach(func(i int, h ParsedHeader) (bool, error) {
			if h.GetName() == "" {
				t.Error("could not find header.name")
				return true, common_errors.KeyNotFoundError{Key: "header.name"}
			}

			if val, err := h.GetValue(ctx); err != nil {
				t.Errorf("could not parse header.valueTmpl - error: %s", err)

				return true, common_errors.KeyNotFoundError{Key: "header.valueTmpl"}
			} else if val == "" {
				t.Error("could not find header.valueTmpl")

				return true, common_errors.KeyNotFoundError{Key: "header.valueTmpl"}
			}

			return false, nil
		})

		if err != nil {
			t.Errorf("could not find Headers. Error - %s", err)
		}
	})
}

func buildContext(integration *Spec) *templates.TemplateContext {
	params := make(map[string]string)
	env := map[string]string{
		"commit_hash":        "hello-commit-hash",
		"commit_author":      "john",
		"hmac_256_signature": "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03",
	}
	ctx := &templates.TemplateContext{Params: params, Env: env}

	integration.Request.Params.ForEach(func(_ int, p *Param) {
		params[p.Name] = "hello-world"
	})

	return ctx
}

func getIntegration(configFile string) (*Spec, error) {
	parser := NewParser()

	return parser.Read(configFile)
}
