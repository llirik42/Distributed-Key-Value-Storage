package grpc

import (
	pb "distributed-key-value-storage/generated/proto"
	"distributed-key-value-storage/src/raft/transport"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type ServerFactory struct{}

func NewServerFactory() *ServerFactory {
	return &ServerFactory{}
}

func (factory ServerFactory) NewServer(
	address string,
	handleRequestForVoteRequest transport.HandleRequestForVoteRequest,
	handleAppendEntriesRequest transport.HandleAppendEntriesRequest,
) (transport.Server, error) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	gRPCServer := grpc.NewServer()

	server := &Server{
		gRPCServer:                  gRPCServer,
		listener:                    listener,
		handleRequestForVoteRequest: handleRequestForVoteRequest,
		handleAppendEntriesRequest:  handleAppendEntriesRequest,
	}

	pb.RegisterRaftServiceServer(gRPCServer, server)
	reflection.Register(gRPCServer)

	return server, nil
}
