package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/set-img/internal/auth"
	"github.com/datewu/set-img/internal/auth/github"
)

func (k *k8sHandler) auth(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		if k.app.Env() == gtea.DevEnv {
			k.user = "datewu"
			next(w, r)
			return
		}
		token, err := handler.GetToken(r, "token")
		if err != nil {
			handler.BadRequestMsg(w, "missing github token query/header/cookie")
			return
		}
		github.DebugGetUser(token)
		time.Sleep(time.Second)
		user, err := github.GetUser(token)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		fmt.Println("user xxxx:", user)
		k.user = user.Login
		k.token = token
		ok, err := auth.Valid(context.Background(), auth.GithubAuth, user.Login, token)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		if !ok {
			handler.InvalidAuthenticationToken(w)
			return
		}
		next(w, r)
	}
	return middle
}
func (m *myHandler) auth(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		if m.app.Env() == gtea.DevEnv {
			m.user = "datewu"
			next(w, r)
			return
		}
		co, err := r.Cookie(github.CookieName)
		if err != nil {
			handler.BadRequestMsg(w, "missing github access_token cookie")
			return
		}
		t := co.Value
		user, err := github.GetUser(t)
		if err != nil {
			handler.ClearSimpleCookie(w, github.CookieName)
			handler.ServerErr(w, err)
			return
		}
		m.user = user.Login
		m.token = t
		next(w, r)
	}
	return middle
}
