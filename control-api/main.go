package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fercevik729/STLKER/control-api/data"
	"github.com/fercevik729/STLKER/control-api/handlers"
	p "github.com/fercevik729/STLKER/watcher-api/protos"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l := log.New(os.Stdout, "control-api", log.LstdFlags)

	// Dial gRPC server
	conn, err := grpc.Dial(":9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Println("[ERROR] dialing gRPC server")
		panic(err)
	}
	defer conn.Close()
	// Create watcher client
	wc := p.NewWatcherClient(conn)
	// Create stock prices database instance
	spdb := data.NewStockPricesDB(wc, l)
	// Create serve mux
	sm := mux.NewRouter()
	// Create handlers
	control := handlers.NewControlHandler(l, spdb)

	// Register handlers
	// TODO: add destination currency parameters
	sm.HandleFunc("/info", control.GetInfo).Queries("ticker", "{ticker:[A-Z]+}")
	sm.HandleFunc("/moreinfo", control.MoreInfo).Queries("ticker", "{ticker:[A-Z]+}")
	sm.HandleFunc("/sub", control.SubscribeTicker).Queries("ticker", "{ticker:[A-Z]+}")
	//sm.Path("/info").Queries("ticker", "{[A-Z]+}", "dest", "{[A-Z]{3}").HandlerFunc(control.GetInfo)
	//sm.Path("/moreinfo").Queries("ticker", "{[A-Z]+}").HandlerFunc(control.GetInfo)
	//sm.Path("/sub").Queries("ticker", "{[A-Z]+}", "dest", "{[A-Z]{3}").HandlerFunc(control.SubscribeTicker)

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
	err = s.Shutdown(tc)
	if err != nil {
		l.Println("Tried shutting down, but got this error:", err)
	}

}
