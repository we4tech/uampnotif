package clients

import "testing"

func TestResponse_IsOK(t *testing.T) {
	t.Run("returns true if valid codes include the returned code", func(t *testing.T) {
		resp := Response{
			Code:       201,
			validCodes: []int{201},
		}

		if !resp.IsOK() {
			t.Error()
		}
	})

	t.Run("returns false if valid codes don't include the returned code", func(t *testing.T) {
		resp := Response{
			Code:       400,
			validCodes: []int{201},
		}

		if resp.IsOK() {
			t.Error()
		}
	})
}
