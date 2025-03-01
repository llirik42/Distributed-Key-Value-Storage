package grpc

import (
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/raft/transport"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientFactory struct{}

func (factory ClientFactory) NewClient(address string) (*transport.Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	gRPCConnection, err := grpc.NewClient(address, opts...)

	if err != nil {
		return nil, errors.Join(errors.New("failed to create gRPC-client"), err)
	}

	client := &Client{
		gRPCClient:     pb.NewRaftServiceClient(gRPCConnection),
		gRPCConnection: gRPCConnection,
	}

	var transportClient transport.Client = client

	return &transportClient, nil
}
