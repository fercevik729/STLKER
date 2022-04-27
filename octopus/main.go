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
	mw "github.com/fercevik729/STLKER/octopus/middleware"
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
	// Register routes
	registerRoutes(sm, control)

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

func registerRoutes(sm *mux.Router, control *handlers.ControlHandler) {
	// Create subrouters and register handlers
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/portfolio/{name}", control.GetPortfolio)
	getR.HandleFunc("/portfolios", control.GetAll)
	getR.HandleFunc("/portfolio/{name}/{ticker}", control.ReadSecurity)
	getR.Use(mw.Authenticate)

	sm.HandleFunc("/info/{ticker}/{currency}", control.GetInfo).Methods("GET")
	sm.HandleFunc("/moreinfo/{ticker}", control.MoreInfo).Methods("GET")

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/portfolio", control.CreatePortfolio)
	postR.HandleFunc("/portfolio/{name}", control.AddSecurity)
	postR.Use(mw.Authenticate)

	// Authentication routes
	sm.HandleFunc("/signup", control.SignUp).Methods("POST")
	sm.HandleFunc("/login", control.LogIn).Methods("POST")
	sm.HandleFunc("/logout", control.LogOut).Methods("GET")
	sm.HandleFunc("/refresh", control.Refresh).Methods("GET")

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/portfolio/{name}/{ticker}/{shares}", control.EditSecurity)
	putR.HandleFunc("/portfolio/{name}", control.UpdatePortfolio)
	putR.Use(mw.Authenticate)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/portfolio/{name}", control.DeletePortfolio)
	deleteR.HandleFunc("/portfolio/{name}/{ticker}", control.DeleteSecurity)
	deleteR.Use(mw.Authenticate)

}
