package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler/sse"
	"github.com/fsnotify/fsnotify"
)

type reloadSSE struct {
	app   *gtea.App
	fs    *fsnotify.Watcher
	dedup time.Duration
}

func (r *reloadSSE) deduplicated() {
	timeout := time.After(r.dedup)
	for {
		select {
		case <-timeout:
			return
		case _, ok := <-r.fs.Events:
			if !ok {
				return
			}
		}
	}
}

func newReloadSSE(app *gtea.App, dirs ...string) *reloadSSE {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	closeFn := func() {
		if err := watcher.Close(); err != nil {
			app.Logger.Err(err)
		}
	}
	app.AddClearFn(closeFn)
	// Add  pathes.
	for _, d := range dirs {
		if err = watcher.Add(d); err != nil {
			app.Logger.Err(err)
			return nil
		}
	}
	return &reloadSSE{
		app:   app,
		fs:    watcher,
		dedup: 300 * time.Millisecond,
	}
}

// Send
func (r *reloadSSE) Send(ctx context.Context, w http.ResponseWriter, f http.Flusher) {
	if r == nil {
		return
	}
	done := ctx.Done()
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
	// defer sse.Shutdown(w, f)

	heartbeat := sse.NewMessage(1)
	r.app.Logger.Info("sse send  heatbeat")
	heartbeat.Send(w, f)
	reload := sse.NewMessage("setTimeout(() => location.reload(), 100)")
	for {
		select {
		case <-done:
			r.app.Logger.Info("sse client disconnected")
			return
		case event, ok := <-r.fs.Events:
			if !ok {
				return
			}
			r.app.Logger.Info("fs event")
			if event.Has(fsnotify.Write) {
				if err := reload.Send(w, f); err != nil {
					r.app.Logger.Err(err)
				}
				r.app.Logger.Info("modified file: " + event.Name)
				r.app.Logger.Info("send reload sse message, and return from for loop")
				r.deduplicated() // not necessary
				return
			}
			r.app.Logger.Info("not write event, continue")
		case err, ok := <-r.fs.Errors:
			if !ok {
				return
			}
			r.app.Logger.Err(err)
			return
		}
	}
}
