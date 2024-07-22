package handler

import "net/http"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func (h *Handler) CallBack(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("code", "200")
}
