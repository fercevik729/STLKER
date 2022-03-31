package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "control-api", log.LstdFlags)

	// Create serve mux
	sm := mux.NewRouter()

	// Create handlers and any middleware

	// CORS for Vuejs UI
	ch := goHandlers.CORS(goHandlers.AllowedOrigins([]string{"https://localhost:3000"}))

	// Create server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      ch(sm),
		ErrorLog:     l,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	go func() {
		l.Println("[DEBUG] Starting control on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Println("[ERROR] starting server", err)
			return
		}
	}()
	// Gracefully stop the server
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	// Blocks until it can read from signal channel
	sig := <-sigChan
	l.Println("Receieved terminate, graceful shutdown", sig)
	tc, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	err := s.Shutdown(tc)
	if err != nil {
		l.Println("Tried shutting down, but got this error:", err)
	}

}
