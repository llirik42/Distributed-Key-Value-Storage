package main

import (
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:9111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Println("listening on :9111")
	}

	grpcServer := grpc.NewServer()
	service := &services.RaftService{}

	pb.RegisterRaftServiceServer(grpcServer, service)
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
