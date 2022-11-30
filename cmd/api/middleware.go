package api

import (
	"context"
	"time"

	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/set-img/internal/auth"
	"github.com/datewu/set-img/internal/author"
)

func checkAuth(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ok, err := auth.Valid(ctx, auth.GithubAuth, token)
	if err != nil {
		jsonlog.Err(err)
		return false, err
	}
	if !ok {
		jsonlog.Err(nil, map[string]any{"token": token})
		return false, nil
	}
	return author.Can(token)
}
