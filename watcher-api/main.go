package main

import (
	"log"
	"net"
	"os"

	"github.com/fercevik729/STLKER/watcher-api/data"
	server "github.com/fercevik729/STLKER/watcher-api/go-server"
	pb "github.com/fercevik729/STLKER/watcher-api/protos"
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
	log.Fatalf("failed to serve: %v\n", gs.Serve(lis))

}
