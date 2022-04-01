package main

import (
	"log"
	"net"
	"os"

	"github.com/fercevik729/STLKER/watcher-api/data"
	"github.com/fercevik729/STLKER/watcher-api/protos"
	"github.com/fercevik729/STLKER/watcher-api/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize the gRPC server instance
	l := log.New(os.Stdout, "watcher-api", log.LstdFlags)
	sp := data.NewStockPrices(l)

	gs := grpc.NewServer()
	cs := server.NewWatcher(sp, l)

	protos.RegisterWatcherServer(gs, cs)

	reflection.Register(gs)

	// Create a tcp socker to listen for incoming connections
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		l.Fatal("Unable to listen, err:", err)
	}
	l.Println("[DEBUG] Listening on port 9092")
	gs.Serve(lis)
}
