package github

import (
	"context"
)

const (
	ghcrRegistry = "https://ghcr.io/v2/"
	username     = "datewu"
)

// Valid ...
func Valide(ctx context.Context, token string) (bool, error) {
	dockerCli := newMicroDockerClient(username, token)
	return dockerCli.CheckToken(ctx, username, token)
}
