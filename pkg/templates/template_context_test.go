package templates

import (
	"testing"
)

func TestFindParam(t *testing.T) {
	ctx := TemplateContext{
		Params: map[string]string{
			"hello": "world",
		},
	}

	t.Run("should find value", func(t *testing.T) {
		if value, _ := ctx.FindParam("hello"); value != "world" {
			t.Errorf("could not find expected value")
		}
	})

	t.Run("should raise error", func(t *testing.T) {
		if _, err := ctx.FindParam("no-key"); err.Error() != "KeyNotFound:no-key" {
			t.Errorf("could not find expected error")
		}
	})
}

func TestParamExists(t *testing.T) {
	ctx := TemplateContext{
		Params: map[string]string{
			"hello": "world",
		},
	}

	t.Run("should find hello param", func(t *testing.T) {
		if !ctx.ParamExists("hello") {
			t.Error("could not find value hello")
		}
	})

	t.Run("should not find not-hello param", func(t *testing.T) {
		if ctx.ParamExists("not-hello") {
			t.Error("could find a value that is not suppose to exist")
		}
	})
}

func TestFindEnv(t *testing.T) {
	ctx := TemplateContext{
		Env: map[string]string{"hello": "world"},
	}

	t.Run("should return value", func(t *testing.T) {
		if value, _ := ctx.FindEnv("hello"); value != "world" {
			t.Error("could not find value for key:hello")
		}
	})

	t.Run(
		"should raise error",
		func(t *testing.T) {
			if _, err := ctx.FindEnv("some-value"); err.Error() != "KeyNotFound:some-value" {
				t.Error("could not find expected error")
			}
		})
}

func TestEnvExists(t *testing.T) {
	ctx := TemplateContext{
		Env: map[string]string{
			"hello": "world",
		},
	}

	t.Run("should find hello env", func(t *testing.T) {
		if !ctx.EnvExists("hello") {
			t.Error("could not find value hello")
		}
	})

	t.Run("should not find not-hello env", func(t *testing.T) {
		if ctx.EnvExists("not-hello") {
			t.Error("could find a value that is not suppose to exist")
		}
	})
}
