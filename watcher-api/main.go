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

	// create gRPC web server with CORS enabled
	grpcWebServer := grpcweb.WrapServer(gs, grpcweb.WithOriginFunc(func(origin string) bool { return true }))

	l.Println("[DEBUG] http server listening at 0.0.0.0:9091")

	srv := &http.Server{
		Handler: grpcWebServer,
		Addr:    "0.0.0.0:9091",
	}
	// Listen on grpc server
	log.Fatal("failed to serve:", srv.ListenAndServe())

}
