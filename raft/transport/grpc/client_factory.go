package grpc

import (
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/raft/transport"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientFactory struct{}

func (factory ClientFactory) NewClient(address string, handleRequestForVoteResponse transport.HandleRequestForVoteResponse, handleAppendEntriesResponse transport.HandleAppendEntriesResponse) (transport.Client, error) {
	//var retryPolicy = `{
	//        "methodConfig": [{
	//    		"name": [{"service": "grpc.examples.echo.Echo"}],
	//            "retryPolicy": {
	//                "MaxAttempts": 20,
	//                "InitialBackoff": ".001s",
	//                "MaxBackoff": ".001s",
	//                "BackoffMultiplier": 1.0,
	//                "RetryableStatusCodes": [ "UNAVAILABLE" ]
	//            }
	//        }]
	//    }`

	// TODO: add retry policy?

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDisableRetry(),
	}

	gRPCConnection, err := grpc.NewClient(address, opts...)

	if err != nil {
		return nil, errors.Join(errors.New("failed to create gRPC-client"), err)
	}

	client := &Client{
		gRPCClient:                   pb.NewRaftServiceClient(gRPCConnection),
		gRPCConnection:               gRPCConnection,
		handleRequestForVoteResponse: handleRequestForVoteResponse,
		handleAppendEntriesResponse:  handleAppendEntriesResponse,
	}

	return client, nil
}
