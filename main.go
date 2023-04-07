package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"

	"log"
	"main/servise"
	"net/http"
	"os"
)

func main() {
	os.Setenv("PORT","8080")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set.")
	}

	listener, err := net.Listen("tcp",fmt.Sprintf(":%s",port))
	if err != nil {
		log.Fatal("port must be setten : ",err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthcheck)
	mux.HandleFunc("/callback", servise.ResponseBot)

	server := &http.Server{
		Handler: mux,
	}

	go server.Serve(listener)

	quit := make(chan os.Signal,1)
	signal.Notify(quit,os.Interrupt)
	<- quit
	log.Printf("stopping server")
	server.Shutdown(context.Background())
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type","application/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}