package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsumeParams(t *testing.T) {
	h := `Bearer realm="https://ghcr.io/token",service="ghcr.io",scope="repository:user/image:pull"`
	// h := `Bearer realm="api.example.com", scope=profile`

	params, v := ConsumeParams(h)
	assert.NotEmpty(t, v)
	assert.NotEmpty(t, params)
}
