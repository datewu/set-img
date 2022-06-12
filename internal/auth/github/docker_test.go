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
