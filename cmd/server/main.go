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

	"github.com/spf13/cobra"

	"github.com/frohwerk/mocktifactory/cmd/server/handlers/api/storage"
	"github.com/frohwerk/mocktifactory/cmd/server/handlers/repository"
)

var (
	path    string
	webhook string
	command = &cobra.Command{Run: start}
)

func init() {
	command.PersistentFlags().StringVarP(&path, "path", "p", "F:/data/.m2/repository", "Path for the storage of artifacts")
	command.PersistentFlags().StringVarP(&webhook, "webhook", "w", "http://localhost:8080/webhook", "URL of the artifacts webhook")
}

func main() {
	command.Use = os.Args[0]
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func start(cmd *cobra.Command, args []string) {
	repo := &repository.Repository{
		Path: path,
		Webhooks: &repository.Webhooks{
			Registered: []repository.Webhook{{Url: webhook}},
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

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() { fmt.Printf("Bye!"); cancel() }()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("ERROR error listening on %v: %v\n", server.Addr, err)
		}
	}()

	<-signals

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("ERROR error during server shutdown: %v\n", err)
	}
}
