package utils

import (
	"testing"
)

func TestConsumeParams(t *testing.T) {
	h := `Bearer realm="https://ghcr.io/token",service="ghcr.io",scope="repository:user/image:pull"`
	// h := `Bearer realm="api.example.com", scope=profile`

	params, v := ConsumeParams(h)
	if len(params) < 1 {
		t.Fatal("params should not be empty")
	}
	if v == "" {
		t.Fatal("newv should not be empty")
	}
}
