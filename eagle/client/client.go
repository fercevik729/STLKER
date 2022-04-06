package main

import (
	"context"
	"log"
	"time"

	pb "github.com/fercevik729/STLKER/eagle/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Couldn't connect to gRPC server", err)
	}
	defer conn.Close()

	client := pb.NewWatcherClient(conn)
	GetInfo(client)
	MoreInfo(client)

}

func GetInfo(client pb.WatcherClient) {
	tr := &pb.TickerRequest{
		Ticker:      "SPY",
		Destination: pb.Currencies_TRY,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tResp, err := client.GetInfo(ctx, tr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", tResp)
}

func MoreInfo(client pb.WatcherClient) {
	tr := &pb.TickerRequest{
		Ticker:      "SPY",
		Destination: pb.Currencies_TRY,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cResp, err := client.MoreInfo(ctx, tr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", cResp)
}
