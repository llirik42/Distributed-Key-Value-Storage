package grpc

import (
	"context"
	pb "distributed-key-value-storage/generated/proto"
	"distributed-key-value-storage/src/raft/dto"
	"distributed-key-value-storage/src/raft/transport"
	"fmt"
	"google.golang.org/grpc"
)

type Client struct {
	index                        int
	address                      string
	gRPCClient                   pb.RaftServiceClient
	gRPCConnection               *grpc.ClientConn
	handleRequestForVoteResponse transport.HandleRequestForVoteResponse
	handleAppendEntriesResponse  transport.HandleAppendEntriesResponse
}

func (client *Client) SendRequestForVote(request dto.RequestVoteRequest) error {
	pbRequest := MapRequestForVoteRequestToGRPC(&request)
	pbResponse, err := client.gRPCClient.RequestForVote(context.Background(), pbRequest)

	if err != nil {
		return fmt.Errorf("failed to send request for vote %s: %w", pbRequest.String(), err)
	}

	response := MapRequestForVoteResponseFromGRPC(pbResponse)
	client.handleRequestForVoteResponse(client, response)

	return nil
}

func (client *Client) SendAppendEntries(request dto.AppendEntriesRequest) error {
	pbRequest := MapAppendEntriesRequestToGRPC(&request)
	pbResponse, err := client.gRPCClient.AppendEntries(context.Background(), pbRequest)

	if err != nil {
		return fmt.Errorf("failed to send append-entries request %s: %w", pbRequest.String(), err)
	}

	response := MapAppendEntriesResponseFromGRPC(pbResponse)
	client.handleAppendEntriesResponse(client, response)

	return nil
}

func (client *Client) GetIndex() int {
	return client.index
}

func (client *Client) GetAddress() string {
	return client.address
}

func (client *Client) Close() error {
	if err := client.gRPCConnection.Close(); err != nil {
		return fmt.Errorf("failed to close gRPC-connection: %w", err)
	}

	return nil
}
