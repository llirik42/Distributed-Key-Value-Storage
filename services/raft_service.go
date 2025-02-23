package services

import (
	"context"
	pb "distributed-algorithms/generated/proto"
)

type RaftService struct {
	pb.RaftServiceServer
}

func (postService *RaftService) RequestForVote(context.Context, *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return &pb.RequestVoteResponse{Term: 1, VoteGranted: false}, nil
}

func (postService *RaftService) AppendEntries(context.Context, *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	return &pb.AppendEntriesResponse{Term: 1, Success: false}, nil
}
