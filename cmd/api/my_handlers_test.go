package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestScaleResourceRejectsInvalidInput(t *testing.T) {
	m := &myHandler{user: "datewu"}
	cases := []url.Values{
		{"ns": {"wu"}, "kind": {"deploy"}, "name": {"web"}, "replicas": {"-1"}},
		{"ns": {"wu"}, "kind": {"deploy"}, "name": {"web"}, "replicas": {"abc"}},
		{"ns": {"wu"}, "kind": {"pod"}, "name": {"web"}, "replicas": {"1"}},
		{"ns": {""}, "kind": {"deploy"}, "name": {"web"}, "replicas": {"1"}},
	}
	for _, form := range cases {
		req := httptest.NewRequest(http.MethodPut, "/my/scale", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		m.scaleResource(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("form %v: got status %d, want %d", form, w.Code, http.StatusBadRequest)
		}
	}
}
