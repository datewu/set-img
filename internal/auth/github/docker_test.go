package github

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectPropertiesHelper(t *testing.T) {
	c := newMicroDockerClient("", "")
	err := c.detectPropertiesHelper(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, c.challenges)
}

func TestParseValueAndParams(t *testing.T) {
	h := `Bearer realm="https://ghcr.io/token",service="ghcr.io",scope="repository:user/image:pull"`
	// h := `Bearer realm="api.example.com", scope=profile`

	params, v := consumeParams(h)
	assert.NotEmpty(t, v)
	assert.NotEmpty(t, params)
}
