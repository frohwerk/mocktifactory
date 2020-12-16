package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frohwerk/mocktifactory/cmd/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/storage/libs-release-local/", handlers.StorageHandler)
	mux.HandleFunc("/libs-release-local/", handlers.DownloadHandler)

	server := http.Server{
		Addr:    ":8091",
		Handler: mux,
	}

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() { fmt.Printf("Bye!"); cancel() }()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("ERROR error listening on %v: %v\n", server.Addr, err)
		}
		done <- true
	}()

	<-signals

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("ERROR error during server shutdown: %v\n", err)
	}
}
