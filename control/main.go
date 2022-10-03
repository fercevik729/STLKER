package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fercevik729/STLKER/control/handlers"
	p "github.com/fercevik729/STLKER/grpc/protos"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-redis/redis/v8"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ring *redis.Ring
var dsn string

func init() {
	// Load environmental variables
	err := godotenv.Load("../vars.env")
	if err != nil {
		panic(errors.New("couldn't load environmental variables from ../vars.env"))
	}
	// Get DB_HOST
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		panic(errors.New("couldn't retrieve DB_HOST"))
	}
	// Get DB_USER
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		panic(errors.New("couldn't retrieve DB_USER"))
	}
	// Get DB_Password
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		panic(errors.New("couldn't retrieve DB_PASSWORD"))
	}
	// Get database name
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic(errors.New("couldn't retrieve DB_NAME"))
	}
	// Get database port
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		panic(errors.New("couldn't retrieve DB_PORT"))
	}
	// Initialize database
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=PDT", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	l := log.New(os.Stdout, "control", log.LstdFlags)

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
	control := handlers.NewControlHandler(l, wc, ring, dsn)
	// Register routes
	registerRoutes(sm, control)
	// CORS for UI (maybe)
	ch := goHandlers.CORS(goHandlers.AllowedOrigins([]string{"*"}))

	// Create server
	s := &http.Server{
		Addr:         "localhost:8080",
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
	sm.Use(control.Logger)
	// Create subrouters and register handlers
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/portfolios/{name}", control.GetPortfolio)
	getR.HandleFunc("/portfolios", control.GetAll)
	getR.HandleFunc("/portfolios/{name}/{ticker}", control.ReadSecurity)
	getR.Use(handlers.Authenticate)

	// Swagger UI
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	docsR := sm.Methods(http.MethodGet).Subrouter()
	docsR.Handle("/docs", sh)
	docsR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Add cache middleware to stockRouter
	stockR := sm.Methods(http.MethodGet).Subrouter()
	stockR.HandleFunc("/stocks/more/{ticker:[A-Z]+}", control.MoreInfo).Methods("GET")
	stockR.HandleFunc("/stocks/{ticker}/{currency}", control.GetInfo).Methods("GET")
	stockR.Use(control.Cache)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/portfolios", control.CreatePortfolio)
	postR.HandleFunc("/portfolios/{name}", control.CreateSecurity)
	postR.Use(handlers.Authenticate)

	// Authentication routes
	sm.HandleFunc("/signup", control.SignUp).Methods("POST")
	sm.HandleFunc("/login", control.LogIn).Methods("POST")
	sm.HandleFunc("/logout", control.LogOut).Methods("POST").Subrouter().Use(handlers.Authenticate)
	sm.HandleFunc("/refresh", control.Refresh).Methods("UPDATE")

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/portfolios/{name}", control.UpdateSecurity)
	putR.HandleFunc("/portfolios", control.UpdatePortfolio)
	putR.Use(handlers.Authenticate)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/portfolios/{name}", control.DeletePortfolio)
	deleteR.HandleFunc("/portfolios/{name}/{ticker}", control.DeleteSecurity)
	deleteR.HandleFunc("/deleteuser", control.DeleteUser)
	deleteR.Use(handlers.Authenticate)

}
