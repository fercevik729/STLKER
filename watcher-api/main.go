package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/fercevik729/STLKER/watcher-api/data"
	pb "github.com/fercevik729/STLKER/watcher-api/protos"
	"github.com/fercevik729/STLKER/watcher-api/server"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize the gRPC server instance
	l := log.New(os.Stdout, "watcher-api", log.LstdFlags)
	sp := data.NewStockPrices(l)
	gs := grpc.NewServer()
	ws := server.NewWatcher(sp, l)
	// Register the Watcher server
	pb.RegisterWatcherServer(gs, ws)
	reflection.Register(gs)

	// Create tcp socket for incoming connections
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		l.Fatal("Unable to listen, err:", err)
	}
	l.Println("[DEBUG] gRPC server listening at localhost:9092")

	go func() {
		// Listen on grpc server
		log.Fatal("failed to serve:", gs.Serve(lis))
	}()

	// Create grpc web server
	grpcWebServer := grpcweb.WrapServer(
		gs, grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	// Listen on web server
	srv := &http.Server{
		Handler: grpcWebServer,
		Addr:    ":9093",
	}
	l.Println("[DEBUG] http server listening at localhost:9093")
	l.Fatalln("failed to serve:", srv.ListenAndServe())
}
