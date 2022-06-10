package api

import (
	"context"
	"net/http"

	"github.com/datewu/gtea/handler"
	"github.com/datewu/set-img/internal/auth"
	"github.com/datewu/set-img/internal/author"
)

func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		h := handler.NewHandleHelper(w, r)
		token, err := handler.GetToken(r, "token")
		if err != nil {
			h.BadRequestErr(err)
			return
		}
		ok, err := auth.Valid(context.Background(), auth.GithubAuth, token)
		if err != nil || !ok {
			h.AuthenticationRequire()
			return
		}
		ok, err = author.Can(token)
		if err != nil {
			h.ServerErr(err)
			return
		}
		if !ok {
			h.NotPermitted()
			return
		}
		next(w, r)
	}
	return http.HandlerFunc(middle)
}
