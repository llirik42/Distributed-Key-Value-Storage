package grpc

import (
	"context"
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/src/raft/domain"
	"distributed-algorithms/src/raft/transport"
	"fmt"
	"google.golang.org/grpc"
)

type Client struct {
	gRPCClient                   pb.RaftServiceClient
	gRPCConnection               *grpc.ClientConn
	handleRequestForVoteResponse transport.HandleRequestForVoteResponse
	handleAppendEntriesResponse  transport.HandleAppendEntriesResponse
}

func (client *Client) SendRequestForVote(request domain.RequestVoteRequest) error {
	pbRequest := &pb.RequestVoteRequest{
		Term:         request.Term,
		CandidateId:  request.CandidateId,
		LastLogIndex: request.LastLogIndex,
		LastLogTerm:  request.LastLogTerm,
	}

	pbResponse, err := client.gRPCClient.RequestForVote(context.Background(), pbRequest)

	if err != nil {
		return fmt.Errorf("failed to send request for vote %s: %w", pbRequest.String(), err)
	}

	response := domain.RequestVoteResponse{
		Term:        pbResponse.Term,
		VoteGranted: pbResponse.VoteGranted,
	}

	client.handleRequestForVoteResponse(&response)

	return nil
}

func (client *Client) SendAppendEntries(request domain.AppendEntriesRequest) error {
	pbRequest := &pb.AppendEntriesRequest{
		Term:         request.Term,
		LeaderId:     request.LeaderId,
		PrevLogIndex: request.PrevLogIndex,
		PrevLogTerm:  request.PrevLogTerm,
		LeaderCommit: request.LeaderCommit,
	}

	pbResponse, err := client.gRPCClient.AppendEntries(context.Background(), pbRequest)

	if err != nil {
		return fmt.Errorf("failed to send append-entries request %s: %w", pbRequest.String(), err)
	}

	response := domain.AppendEntriesResponse{
		Term:    pbResponse.Term,
		Success: pbResponse.Success,
	}

	client.handleAppendEntriesResponse(&response)

	return nil
}

func (client *Client) Close() error {
	if err := client.gRPCConnection.Close(); err != nil {
		return fmt.Errorf("failed to close gRPC-connection: %w", err)
	}

	return nil
}
