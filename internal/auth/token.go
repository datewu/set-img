package auth

import (
	"context"
	"errors"

	"github.com/datewu/set-img/internal/auth/github"
)

type authType int

const (
	// GithubAuth ...
	GithubAuth authType = iota
)

// ErrInvalidToken ...
var ErrInvalidToken = errors.New("invalid token")

// Valid ...
func Valid(ctx context.Context, kind authType, username, token string) (bool, error) {
	switch kind {
	case GithubAuth:
		return github.Valide(ctx, username, token)
	}
	return false, ErrInvalidToken
}
