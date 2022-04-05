package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/fercevik729/STLKER/watcher-api/data"
	server "github.com/fercevik729/STLKER/watcher-api/go-server"
	pb "github.com/fercevik729/STLKER/watcher-api/protos"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcMux struct {
	*grpcweb.WrappedGrpcServer
}

// Handler routes requests to grpc or the regular http server
func (m *grpcMux) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.IsGrpcWebRequest(r) {
			m.ServeHTTP(w, r)
			return
		} else if m.IsAcceptableGrpcCorsRequest(r) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		//next.ServeHTTP(w, r)
	})
}
func main() {
	l := log.New(os.Stdout, "watcher-api", log.LstdFlags)

	// Create tcp socket for incoming connections
	lis, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		l.Fatal("Unable to listen, err:", err)
	}
	l.Println("[DEBUG] gRPC server listening at 0.0.0.0:9090")

	// Initialize the gRPC server instance
	sp := data.NewStockPrices(l)
	gs := grpc.NewServer()
	ws := server.NewWatcher(sp, l)
	// Register the Watcher server
	pb.RegisterWatcherServer(gs, ws)
	reflection.Register(gs)

	// Start up gRPC server
	go func() {
		log.Fatalf("failed to serve: %v\n", gs.Serve(lis))
	}()

	// Create gRPC web server with CORS enabled
	grpcWebServer := grpcweb.WrapServer(gs,
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return origin == "http://localhost:8080" || origin == "localhost:9090" || origin == "localhost:3000"
		}))
	mux := grpcMux{grpcWebServer}

	// Create http router and fileserver
	r := http.NewServeMux()
	webapp := http.FileServer(http.Dir("ui/stlkerapp/build"))

	// Register handlers
	r.Handle("/", mux.Handler(webapp))

	// Start up http server
	l.Println("[DEBUG] http server listening at 0.0.0.0:8080")
	srv := &http.Server{
		Handler: grpcWebServer,
		Addr:    "0.0.0.0:8080",
	}
	// Listen on grpc server
	log.Fatal("failed to serve:", srv.ListenAndServe())

}
