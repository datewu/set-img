package api

import (
	"net/http"

	"github.com/datewu/set-img/internal/auth"
	"github.com/datewu/set-img/internal/author"
	"github.com/datewu/toushi"
)

func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		token, err := toushi.GetToken(r, "token")
		if err != nil {
			toushi.HandleBadRequestErr(err)(w, r)
			return
		}
		ok, err := auth.Valid(token)
		if err != nil || !ok {
			toushi.HandleAuthenticationRequire(w, r)
			return
		}
		ok, err = author.Can(token)
		if err != nil {
			toushi.HandleServerErr(err)(w, r)
			return
		}
		if !ok {
			toushi.HandleNotPermitted(w, r)
			return
		}
		next(w, r)
	}
	return http.HandlerFunc(middle)
}
