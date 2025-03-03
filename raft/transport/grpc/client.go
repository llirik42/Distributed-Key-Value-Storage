package grpc

import (
	"context"
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/raft/domain"
	"errors"
	"google.golang.org/grpc"
)

type Client struct {
	gRPCClient     pb.RaftServiceClient
	gRPCConnection *grpc.ClientConn
}

func (client *Client) SendRequestForVote(request domain.RequestVoteRequest) (*domain.RequestVoteResponse, error) {
	pbRequest := &pb.RequestVoteRequest{
		Term:         request.Term,
		CandidateId:  request.CandidateId,
		LastLogIndex: request.LastLogIndex,
		LastLogTerm:  request.LastLogTerm,
	}

	pbResponse, pbErr := client.gRPCClient.RequestForVote(context.Background(), pbRequest)

	if pbErr != nil {
		return nil, errors.Join(errors.New("failed to send request for vote: "+pbRequest.String()), pbErr)
	}

	response := domain.RequestVoteResponse{
		Term:        pbResponse.Term,
		VoteGranted: pbResponse.VoteGranted,
	}

	return &response, nil
}

func (client *Client) SendAppendEntries(request domain.AppendEntriesRequest) (*domain.AppendEntriesResponse, error) {
	// TODO: add retry policy?

	pbRequest := &pb.AppendEntriesRequest{
		Term:         request.Term,
		LeaderId:     request.LeaderId,
		PrevLogIndex: request.PrevLogIndex,
		PrevLogTerm:  request.PrevLogTerm,
		LeaderCommit: request.LeaderCommit,
	}

	pbResponse, pbErr := client.gRPCClient.AppendEntries(context.Background(), pbRequest)

	if pbErr != nil {
		return nil, errors.Join(errors.New("failed to send append entries: "+pbRequest.String()), pbErr)
	}

	response := domain.AppendEntriesResponse{
		Term:    pbResponse.Term,
		Success: pbResponse.Success,
	}

	return &response, nil
}

// Close TODO: is this method ever used?
func (client *Client) Close() error {
	err := client.gRPCConnection.Close()

	if err != nil {
		return errors.Join(errors.New("failed to close gRPC-client"), err)
	}

	return nil
}
