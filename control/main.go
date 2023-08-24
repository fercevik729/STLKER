package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"io"
	"log"
	"log/slog"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: change localhost to names of containers when dockerizing
var ring *redis.Ring
var dsn string

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	// Colorize the level output
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	// Gather the fields
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	//Marshal the fields for JSON output
	b, err := json.MarshalIndent(fields, "", " ")
	if err != nil {
		return err
	}
	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

// NewPrettyHandler constructs a new PrettyHandler struct and returns a pointer to it
func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "control", 0),
	}
	return h
}

func init() {
	// Get DB_Password
	err := godotenv.Load("../app.env")
	if err != nil {
		log.Println("[WARN] Couldn't load env vars from ../app.env")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	// Panic if it can't
	if dbPassword == "" {
		log.Fatalln(errors.New("couldn't get DB_PASSWORD"))
	}
	// Initialize database
	dsn = fmt.Sprintf("host=localhost user=furkanercevik password=%v dbname=stlker port=5432 sslmode=disable",
		dbPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(errors.New("couldn't create connection to database"))
	}
	// Migrate the schemas
	err = db.AutoMigrate(&handlers.Portfolio{}, &handlers.Security{}, &handlers.User{})
	if err != nil {
		log.Fatalln(errors.New("couldn't migrate schemas"))
	}

	// Initialize redis options
	ring = redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "localhost:6379",
		},
	})

}

func main() {
	// Init zap logger
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := NewPrettyHandler(os.Stderr, opts)
	l := slog.New(handler)

	// Dial gRPC server
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error(fmt.Sprintf("Couldn't dial gRPC server: %s", err))
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(conn)
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
		Addr:         "0.0.0.0:8080",
		Handler:      ch(sm),
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	go func() {
		l.Debug("Starting control on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("unable to start server", err)
			return
		}
	}()
	// Gracefully stop the server
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	// Blocks until it can read from signal channel
	sig := <-sigChan
	l.Info("Received terminate, graceful shutdown", "sig", sig)
	tc, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	err = s.Shutdown(tc)
	if err != nil {
		l.Warn("Tried shutting down, but got this", "err", err)
	}

}

func registerRoutes(sm *mux.Router, control *handlers.ControlHandler) {
	sm.Use(control.Logger)
	// Create sub-routers and register handlers
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
