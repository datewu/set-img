package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler/sse"
)

type reloadSSE struct {
	app *gtea.App
}

// Pour
func (r reloadSSE) Pour(w http.ResponseWriter, f http.Flusher) {
	if r.app.Env() != gtea.DevEnv {
		e := sse.NewEvent("reload", "not in development mode. bye!")
		if err := e.Send(w, f); err != nil {
			r.app.Logger.Err(err)
		}
		if err := sse.Shutdown(w, f); err != nil {
			r.app.Logger.Err(err)
		}
		return
	}
	e := sse.NewEvent("reload", "hello world")
	if err := e.Send(w, f); err != nil {
		r.app.Logger.Err(err)
	}
	// time.Sleep(3 * time.Second)
	// e = sse.NewMessage("location.reload()")
	// if err := e.Send(w, f); err != nil {
	// 	r.app.Logger.Err(err)
	// }
	if err := sse.Shutdown(w, f); err != nil {
		r.app.Logger.Err(err)
	}
}
