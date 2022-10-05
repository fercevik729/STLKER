package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/fercevik729/STLKER/grpc/data"
	pb "github.com/fercevik729/STLKER/grpc/protos"
	server "github.com/fercevik729/STLKER/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var apiKey string

func init() {
	// Load the API_KEY before hand
	apiKey = os.Getenv("API_KEY")

	// Panic if it can't
	if apiKey == "" {
		panic(errors.New("couldn't get API_KEY"))
	}
	log.Println("[INFO] API_KEY is:", apiKey)

}

func main() {
	l := log.New(os.Stdout, "grpc", log.LstdFlags)
	// Create tcp socket for incoming connections
	lis, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		l.Fatal("Unable to listen, err:", err)
	}
	l.Println("[DEBUG] Starting grpc on port 9090")

	// Initialize the gRPC server instance
	sp := data.NewStockPrices(l)
	gs := grpc.NewServer()
	ws := server.NewWatcher(sp, l, apiKey)
	// Register the Watcher server
	pb.RegisterWatcherServer(gs, ws)
	reflection.Register(gs)

	// Start up gRPC server
	log.Fatalf("failed to serve: %v\n", gs.Serve(lis))

}
