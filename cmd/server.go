package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"

	"github.com/idirall22/gateway/pb"
)

type server struct{}

func (s *server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Msg: "Received message: " + req.Msg}, nil
}

func (s *server) EchoStream(req *pb.EchoRequest, stream pb.EchoService_EchoStreamServer) error {
	for i := 0; i < 5; i++ {
		err := stream.Send(&pb.EchoResponse{Msg: "From server: " + req.Msg + fmt.Sprintf(" %d", i)})
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func grpcServer(l net.Listener, cred credentials.TransportCredentials) error {
	s := grpc.NewServer(
	// grpc.Creds(cred),
	)
	pb.RegisterEchoServiceServer(s, &server{})
	return s.Serve(l)
}

func restServer(
	l net.Listener,
	echoServer pb.EchoServiceServer,
	endpoint string,
) error {
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatal(err)
	}
	return http.Serve(l, mux)
	// return http.ServeTLS(l, mux, "cert/server-cert.pem", "cert/server-key.pem")
}

func main() {
	svr := os.Getenv("SERVER_TYPE")
	if svr == "" {
		svr = "grpc"
	}

	svrPort := os.Getenv("SERVER_PORT")
	if svrPort == "" {
		svrPort = ":8080"
	}

	svrGrpcPort := os.Getenv("SERVER_GRPC_PORT")
	if svrPort == "" {
		svrPort = ":80"
	}

	endpoint := os.Getenv("SERVER_GRPC_ENDPOINT")
	if endpoint == "" && svr == "rest" {
		endpoint = "grpc-server-service.default.svc.cluster.local" + svrGrpcPort
	}

	l, err := net.Listen("tcp", svrPort)
	if err != nil {
		log.Fatal(err)
	}

	// serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	// if err != nil {
	// 	log.Fatal("Tls error: ", err)
	// }
	config := &tls.Config{
		// Certificates: []tls.Certificate{serverCert},
		ClientAuth: tls.NoClientCert,
	}
	cred := credentials.NewTLS(config)

	if svr == "grpc" {
		log.Println("Start grpc server on", svrPort)
		log.Fatal(grpcServer(l, cred))
	} else {
		log.Printf("Start rest server on %s and grpc on %s", svrPort, endpoint)
		log.Fatal(restServer(l, &server{}, endpoint))
	}
}
