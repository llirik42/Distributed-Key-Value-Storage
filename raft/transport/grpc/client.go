package grpc

import (
	"context"
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/transport"
	"errors"
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

	pbResponse, pbErr := client.gRPCClient.RequestForVote(context.Background(), pbRequest)

	if pbErr != nil {
		return errors.Join(errors.New("failed to send request for vote: "+pbRequest.String()), pbErr)
	}

	response := domain.RequestVoteResponse{
		Term:        pbResponse.Term,
		VoteGranted: pbResponse.VoteGranted,
	}

	client.handleRequestForVoteResponse(&response)

	return nil
}

func (client *Client) SendAppendEntries(request domain.AppendEntriesRequest) error {
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
		return errors.Join(errors.New("failed to send append entries: "+pbRequest.String()), pbErr)
	}

	response := domain.AppendEntriesResponse{
		Term:    pbResponse.Term,
		Success: pbResponse.Success,
	}

	client.handleAppendEntriesResponse(&response)

	return nil
}

func (client *Client) Close() error {
	err := client.gRPCConnection.Close()

	if err != nil {
		return errors.Join(errors.New("failed to close gRPC-client"), err)
	}

	return nil
}
