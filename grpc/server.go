package grpc

import (
	"context"
	pb "distributed-algorithms/generated/proto"
	"distributed-algorithms/raft"
)

type RaftServiceServer struct {
	pb.RaftServiceServer

	Node *raft.Node
}

func (raftService *RaftServiceServer) RequestForVote(ctx context.Context, request *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
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

func (raftService *RaftServiceServer) AppendEntries(ctx context.Context, request *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	// TODO: add checks related to log entries

	currentTerm := raftService.Node.GetCurrentTerm()
	success := request.Term >= currentTerm
	return &pb.AppendEntriesResponse{Term: currentTerm, Success: success}, nil
}
