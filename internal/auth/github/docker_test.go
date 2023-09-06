package github

import (
	"context"
	"testing"
)

func TestDetectPropertiesHelper(t *testing.T) {
	c := newMicroDockerClient("", "")
	err := c.detectPropertiesHelper(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(c.challenges) < 1 {
		t.Fatal("should not be empty")
	}
}
