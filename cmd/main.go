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

	"github.com/frohwerk/mocktifactory/cmd/handlers/api/storage"
	"github.com/frohwerk/mocktifactory/cmd/handlers/repository"
)

func main() {
	repo := &repository.Repository{
		Path: "F:/data/.m2/repository",
		Webhooks: &repository.Webhooks{
			Registered: []repository.Webhook{{Url: "http://localhost:8080/webhook"}},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/storage/libs-release-local/", storage.CreateHandler(repo))
	mux.Handle("/libs-release-local/", repo)

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
