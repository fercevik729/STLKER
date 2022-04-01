package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/fercevik729/STLKER/watcher-api/data"
	"github.com/fercevik729/STLKER/watcher-api/protos"
	"github.com/fercevik729/STLKER/watcher-api/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func testing() {
	// TODO: fix Stock struct's struct tags for JSON response body
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + "SPY" + "&apikey=" + "ECI3WFKP22I15OFC"
	resp, _ := http.Get(url)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(b)

}

func main() {
	testing()
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
