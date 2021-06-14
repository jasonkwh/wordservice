package main

import (
	"flag"
	"log"
	"net"
	"google.golang.org/grpc"
	"deltatre_grpc/proto/wordservice"
	"deltatre_grpc/service"
)

var grpcEndpoint = flag.String("grpc-endpoint", "localhost:9000", "gRPC server endpoint")

func main()  {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9000")
	if err != nil {
		return err
	}

	s := service.Server{}
	s.SetDefaultWords()
	grpcServer := grpc.NewServer()

	wordservice.RegisterWordServiceServer(
		grpcServer,
		&s,
	)

	log.Printf("gRPC server listening on: %s", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}