package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/datewu/set-img/auth"
	"github.com/datewu/set-img/author"
	"github.com/datewu/toushi"
)

func checkAuth(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			toushi.BadRequestResponse(err)(w, r)
			return
		}
		ok, err := auth.Valid(token)
		if err != nil || !ok {
			toushi.AuthenticationRequireResponse(w, r)
			return
		}
		ok, err = author.Can(token)
		if err != nil {
			toushi.ServerErrResponse(err)(w, r)
			return
		}
		if !ok {
			toushi.NotPermittedResponse(w, r)
			return
		}
		next(w, r)
	}
	return http.HandlerFunc(middle)
}

func extractToken(r *http.Request) (string, error) {
	q := toushi.ReadString(r.URL.Query(), "token", "")
	if q != "" {
		return q, nil
	}
	authorizationHeader := r.Header.Get("Authorization")

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("bad authorization header")
	}
	return headerParts[1], nil
}
