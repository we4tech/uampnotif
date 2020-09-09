package integrations

import "testing"

func TestIsEmpty(t *testing.T) {
	t.Run("should return true when param is not empty", func(t *testing.T) {
		ps := Params{
			Param{},
		}

		if ps.IsEmpty() {
			t.Error("could not find IsEmpty = false")
		}
	})

	t.Run("should return false when Params is empty", func(t *testing.T) {
		ps := Params{}

		if !ps.IsEmpty() {
			t.Error("could not find IsEmpty = true")
		}
	})
}
