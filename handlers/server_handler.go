package handlers

import (
	"distributed-algorithms/dto"
	pb "distributed-algorithms/generated/proto"
)

type ServerHandler struct {
}

func (handler *ServerHandler) RequestForVote(request dto.RequestVoteRequest) (dto.RequestVoteRequest, error) {
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

func (handler *ServerHandler) AppendEntries(ctx context.Context, request *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	// TODO: add checks related to log entries

	currentTerm := raftService.Node.GetCurrentTerm()
	success := request.Term >= currentTerm
	return &pb.AppendEntriesResponse{Term: currentTerm, Success: success}, nil
}
