package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fercevik729/STLKER/octopus/handlers"
	p "github.com/fercevik729/STLKER/watcher-api/protos"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l := log.New(os.Stdout, "control-api", log.LstdFlags)

	// Dial gRPC server
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Println("[ERROR] dialing gRPC server")
		panic(err)
	}
	defer conn.Close()
	// Create watcher client
	wc := p.NewWatcherClient(conn)
	// Create serve mux
	sm := mux.NewRouter()
	// Create handlers
	control := handlers.NewControlHandler(l, wc)

	// Create subrouters and register handlers
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/portfolio/{name}", control.GetPortfolio)
	getRouter.HandleFunc("/info/{ticker}/{currency}", control.GetInfo)
	getRouter.HandleFunc("/moreinfo/{ticker}", control.MoreInfo)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/portfolio", control.CreatePortfolio)
	// Change to use json
	postRouter.HandleFunc("/portfolio/{name}/{ticker}/{shares}", control.AddSecurity)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/portfolio/{name}/{ticker}/{shares}", control.EditSecurity)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/portfolio/{name}", control.DeletePortfolio)
	deleteRouter.HandleFunc("/portfolio/{name}/{ticker}", control.DeleteSecurity)

	// CORS for UI
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
	err = s.Shutdown(tc)
	if err != nil {
		l.Println("Tried shutting down, but got this error:", err)
	}

}
