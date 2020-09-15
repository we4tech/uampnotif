package notifications

import (
	"github.com/we4tech/uampnotif/pkg/templates"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	t.Run("should return true when empty", func(t *testing.T) {
		params := &Params{"key": "value"}

		if params.IsEmpty() {
			t.Error()
		}
	})

	t.Run("should return false when not empty", func(t *testing.T) {
		params := Params{}

		if !params.IsEmpty() {
			t.Error()
		}
	})
}

func TestGetValue(t *testing.T) {
	t.Run("should render value", func(t *testing.T) {
		params := &Params{
			"key1": "hello {{.FindParam \"name\"}}",
		}

		ctx := &templates.TemplateContext{
			Params: map[string]string{"name": "uampnotif"}}

		if value, err := params.GetValue(ctx, "key1"); err != nil {
			t.Fatalf("could not parse value. error - %s", err)
		} else if value != "hello uampnotif" {
			t.Fail()
		}
	})

	t.Run("should raise rendering error", func(t *testing.T) {
		params := &Params{
			"key1": "hello {{.FindParam \"name\"}}",
		}

		ctx := &templates.TemplateContext{}

		if _, err := params.GetValue(ctx, "key1"); err == nil {
			t.Fatal("did not raise error")
		}
	})
}
