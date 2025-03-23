package grpc

import (
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/src/raft/transport"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ClientFactory struct{}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{}
}

func (factory ClientFactory) NewClient(
	index int,
	address string,
	handleRequestForVoteResponse transport.HandleRequestForVoteResponse,
	handleAppendEntriesResponse transport.HandleAppendEntriesResponse,
) (transport.Client, error) {
	// TODO: connect params should depend on RAFT-timeouts
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  10 * time.Millisecond,
				Multiplier: 1,
				Jitter:     1,
				MaxDelay:   10 * time.Millisecond,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
	}

	gRPCConnection, err := grpc.NewClient(address, opts...)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC-connection with %s: %w", address, err)
	}

	client := &Client{
		index:                        index,
		gRPCClient:                   pb.NewRaftServiceClient(gRPCConnection),
		gRPCConnection:               gRPCConnection,
		handleRequestForVoteResponse: handleRequestForVoteResponse,
		handleAppendEntriesResponse:  handleAppendEntriesResponse,
	}

	return client, nil
}
