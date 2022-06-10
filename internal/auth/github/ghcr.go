package github

import (
	"context"

	"github.com/rs/zerolog/log"
)

const (
	ghcrRegistry = "https://ghcr.io/v2/"
	username     = "datewu"
)

// Valid ...
func Valide(ctx context.Context, token string) (bool, error) {
	dockerCli := newMicroDockerClient(username, token)
	log.Info().Str("token", token).Str("username", username).
		Msg("check token")
	ok, err := dockerCli.CheckToken(ctx, username, token)
	if err != nil {
		log.Err(err).Msg("check token failed")
		return false, err
	}
	return ok, nil
}
