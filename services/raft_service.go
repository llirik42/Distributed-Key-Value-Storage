package services

import (
	"context"
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/raft"
)

type RaftService struct {
	pb.RaftServiceServer
	Node *raft.Node
}

func (raftService *RaftService) RequestForVote(ctx context.Context, request *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	// TODO: add checks about candidate's log

	currentTerm := raftService.Node.GetCurrentTerm()

	var voteGranted bool

	if request.Term < currentTerm {
		voteGranted = false
	} else {
		voteGranted = raftService.Node.Vote(request.CandidateId)
	}

	return &pb.RequestVoteResponse{Term: currentTerm, VoteGranted: voteGranted}, nil
}

func (raftService *RaftService) AppendEntries(ctx context.Context, request *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	currentTerm := raftService.Node.GetCurrentTerm()
	success := request.Term >= currentTerm // TODO: add checks related to log entries
	return &pb.AppendEntriesResponse{Term: currentTerm, Success: success}, nil
}
