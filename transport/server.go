package transport

import (
	"context"
	"distributed-algorithms/dto"
	pb "distributed-algorithms/generated/proto"
)

type RaftServiceServer struct {
	pb.RaftServiceServer
	HandleRequestForVoteRequest func(request dto.RequestVoteRequest) (dto.RequestVoteResponse, error)
	HandleAppendEntriesRequest  func(request dto.AppendEntriesRequest) (dto.AppendEntriesResponse, error)
}

func (server *RaftServiceServer) RequestForVote(ctx context.Context, request *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	result, err := server.HandleRequestForVoteRequest(dto.RequestVoteRequest{
		Term:         request.Term,
		CandidateId:  request.CandidateId,
		LastLogIndex: request.LastLogIndex,
		LastLogTerm:  request.LastLogTerm,
	})

	if err != nil {
		return nil, err
	}

	return &pb.RequestVoteResponse{Term: result.Term, VoteGranted: result.VoteGranted}, nil
}

func (server *RaftServiceServer) AppendEntries(ctx context.Context, request *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	result, err := server.HandleAppendEntriesRequest(dto.AppendEntriesRequest{
		Term:         request.Term,
		LeaderId:     request.LeaderId,
		PrevLogIndex: request.PrevLogIndex,
		PrevLogTerm:  request.PrevLogTerm,
		LeaderCommit: request.LeaderCommit,
	})

	if err != nil {
		return nil, err
	}

	return &pb.AppendEntriesResponse{Term: result.Term, Success: result.Success}, nil
}
