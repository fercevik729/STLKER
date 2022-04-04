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
	// Initialize the gRPC server instance
	l := log.New(os.Stdout, "watcher-api", log.LstdFlags)
	sp := data.NewStockPrices(l)
	gs := grpc.NewServer()
	ws := server.NewWatcher(sp, l)
	// Register the Watcher server
	pb.RegisterWatcherServer(gs, ws)
	reflection.Register(gs)

	// Create tcp socket for incoming connections
	lis, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		l.Fatal("Unable to listen, err:", err)
	}
	l.Println("[DEBUG] gRPC server listening at 0.0.0.0:9090")

	// Listen on grpc server
	log.Fatal("failed to serve:", gs.Serve(lis))

}
