package github

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	ghcrRegistry = "https://ghcr.io/v2/"
	username     = "datewu"
)

// Valid ...
func Valide(ctx context.Context, token string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ghcrRegistry, nil)
	if err != nil {
		return false, err
	}
	log.Info().Str("token", token).
		Msg("going to validate github token")
	req.SetBasicAuth(username, token)
	cli := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Info().Msgf("response body: %s", body)
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}
