package main

import (
	"io"
	"log"

	"golang.org/x/net/context"

	"github.com/idirall22/gateway/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("172.17.0.2:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)
	// res, err := client.Echo(context.Background(), &pb.EchoRequest{Msg: "Hello server"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	stream, err := client.EchoStream(context.Background(), &pb.EchoRequest{Msg: "Idir"})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("No more data")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res.Msg)
	}
}
