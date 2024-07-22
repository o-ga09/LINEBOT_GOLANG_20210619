package main

import (
	"main/handler"
	"net/http"

	"github.com/syumai/workers"
)

func main() {
	handler := handler.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Healthcheck)
	mux.HandleFunc("/callback", handler.CallBack)

	workers.Serve(mux)
}
