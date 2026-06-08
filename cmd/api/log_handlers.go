package api

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/handler/sse"
	"github.com/datewu/set-img/internal/k8s"
)

type logSSE struct {
	ns        string
	kind      string
	name      string
	container string
	pod       string
}

func (l *logSSE) Send(ctx context.Context, w http.ResponseWriter, f http.Flusher) {
	tailLines := int64(100)
	rc, err := k8s.StreamPodLogs(ctx, l.ns, l.pod, l.container, tailLines)
	if err != nil {
		errMsg := sse.NewEvent("error", fmt.Sprintf("Failed to stream logs: %v", err))
		errMsg.Send(w, f)
		sse.Shutdown(w, f)
		return
	}
	defer rc.Close()

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			line := scanner.Text()
			msg := sse.NewMessage(line)
			if err := msg.Send(w, f); err != nil {
				return
			}
		}
	}
	if err := scanner.Err(); err != nil {
		select {
		case <-ctx.Done():
			// client disconnected, ignore error
		default:
			errMsg := sse.NewEvent("error", fmt.Sprintf("Log stream error: %v", err))
			errMsg.Send(w, f)
		}
	}
	sse.Shutdown(w, f)
}

func (m *myHandler) logs(w http.ResponseWriter, r *http.Request) {
	ns := handler.ReadQuery(r, "ns", "")
	kind := handler.ReadQuery(r, "kind", "")
	name := handler.ReadQuery(r, "name", "")
	container := handler.ReadQuery(r, "container", "")
	pod := handler.ReadQuery(r, "pod", "")

	if ns == "" || kind == "" || name == "" || container == "" {
		handler.BadRequestMsg(w, "missing required query params: ns, kind, name, container")
		return
	}

	// If no specific pod was requested, pick the first pod
	if pod == "" {
		pods, err := k8s.ListWorkloadPods(r.Context(), ns, kind, name)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		if len(pods) == 0 {
			handler.BadRequestMsg(w, "no pods found for this workload")
			return
		}
		pod = pods[0]
	}

	streamer := &logSSE{
		ns:        ns,
		kind:      kind,
		name:      name,
		container: container,
		pod:       pod,
	}
	sse.SSE(w, r, streamer)
}

func (m *myHandler) listPods(w http.ResponseWriter, r *http.Request) {
	ns := handler.ReadQuery(r, "ns", "")
	kind := handler.ReadQuery(r, "kind", "")
	name := handler.ReadQuery(r, "name", "")

	if ns == "" || kind == "" || name == "" {
		handler.BadRequestMsg(w, "missing required query params: ns, kind, name")
		return
	}

	pods, err := k8s.ListWorkloadPods(r.Context(), ns, kind, name)
	if err != nil {
		handler.ServerErr(w, err)
		return
	}

	// Render pods as HTML options for an HTMX pod selector
	html := ""
	for _, p := range pods {
		html += fmt.Sprintf(`<option value="%s">%s</option>`, p, p)
	}
	if html == "" {
		html = `<option value="">No pods found</option>`
	}
	handler.OKText(w, html)
}
