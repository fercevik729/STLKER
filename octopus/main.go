package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	p "github.com/fercevik729/STLKER/eagle/protos"
	"github.com/fercevik729/STLKER/octopus/handlers"
	"github.com/go-redis/redis/v8"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ring *redis.Ring

// TODO: create swagger documentation
func init() {
	// Get database name
	dbName, err := handlers.ReadEnvVar("DB_NAME")
	if err != nil {
		panic(err)
	}
	// Initialize database
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Migrate the schemas
	db.AutoMigrate(&handlers.Portfolio{}, &handlers.Security{}, &handlers.User{})

	// Initialize redis options
	ring = redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
		},
	})

}

func main() {
	l := log.New(os.Stdout, "octopus", log.LstdFlags)

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
	control := handlers.NewControlHandler(l, wc, ring)
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
		l.Println("[DEBUG] Starting octopus on port 8080")

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
	getR.HandleFunc("/portfolios/{name}", control.GetPortfolio)
	getR.HandleFunc("/portfolios", control.GetAll)
	getR.HandleFunc("/portfolios/{name}/{ticker}", control.ReadSecurity)
	// Add cache middleware to getRouter
	getR.Use(handlers.Authenticate, control.Cache)

	sm.HandleFunc("/info/{ticker}/{currency}", control.GetInfo).Methods("GET")
	sm.HandleFunc("/moreinfo/{ticker}", control.MoreInfo).Methods("GET")

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/portfolios", control.CreatePortfolio)
	postR.HandleFunc("/portfolios/{name}", control.CreateSecurity)
	postR.Use(handlers.Authenticate)

	// Authentication routes
	sm.HandleFunc("/signup", control.SignUp).Methods("POST")
	sm.HandleFunc("/login", control.LogIn).Methods("POST")
	sm.HandleFunc("/logout", control.LogOut).Methods("GET")
	sm.HandleFunc("/refresh", control.Refresh).Methods("GET")

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/portfolios/{name}", control.UpdateSecurity)
	putR.HandleFunc("/portfolios", control.UpdatePortfolio)
	putR.Use(handlers.Authenticate)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/portfolios/{name}", control.DeletePortfolio)
	deleteR.HandleFunc("/portfolios/{name}/{ticker}", control.DeleteSecurity)
	deleteR.Use(handlers.Authenticate)

}
