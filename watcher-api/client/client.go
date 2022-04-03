package main

import (
	"context"
	"io"
	"log"

	pb "github.com/fercevik729/STLKER/watcher-api/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial(":9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Couldn't connect to gRPC server", err)
	}
	defer conn.Close()

	client := pb.NewWatcherClient(conn)
	SubscribeTicker(client)

}

func SubscribeTicker(client pb.WatcherClient) {
	tr := &pb.TickerRequest{
		Ticker:      "SPY",
		Destination: pb.Currencies_TRY,
	}

	stream, err := client.SubscribeTicker(context.Background(), tr)
	if err != nil {
		log.Fatalln("Error subscribing ticker:", err)
	}

	for {
		price, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error getting price response:", err)
		}

		log.Printf("%#v", price)
	}
	stream.CloseSend()
}
