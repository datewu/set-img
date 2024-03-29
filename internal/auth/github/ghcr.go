package github

import (
	"context"

	"github.com/datewu/gtea/jsonlog"
)

const (
	ghcrRegistry = "https://ghcr.io/v2/"
)

// Valid ...
func Valide(ctx context.Context, username, token string) (bool, error) {
	dockerCli := newMicroDockerClient(username, token)
	jsonlog.Info("check token", map[string]any{"token": token, "username": username})
	ok, err := dockerCli.CheckToken(ctx, username, token)
	if err != nil {
		jsonlog.Err(err, map[string]any{"token": token,
			"username": username, "msg": "check token failed"})
		return false, err
	}
	return ok, nil
}
