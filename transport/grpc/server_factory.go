package grpc

import (
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type ServerFactory struct{}

func (factory *ServerFactory) NewServer(address string, handleRequestForVoteRequest transport.HandleRequestForVoteRequest, handleAppendEntriesRequest transport.HandleAppendEntriesRequest) (*Server, error) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
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
