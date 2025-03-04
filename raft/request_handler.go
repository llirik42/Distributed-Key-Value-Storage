package raft

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/utils"
	"log"
)

type RequestHandler struct {
	ctx *context.Context
}

func NewRequestHandler(ctx *context.Context) *RequestHandler {
	return &RequestHandler{ctx: ctx}
}

func (handler *RequestHandler) HandleRequestVoteRequest(request domain.RequestVoteRequest) (*domain.RequestVoteResponse, error) {
	// TODO: add checks about candidate's log

	ctx := handler.ctx

	log.Printf("Node \"%s\" received request of vote: %v", ctx.GetNodeId(), request)

	currentTerm := ctx.GetCurrentTerm()
	var voteGranted bool

	if request.Term < currentTerm {
		voteGranted = false
	} else {
		voteGranted = ctx.Vote(request.CandidateId)
	}

	return &domain.RequestVoteResponse{Term: currentTerm, VoteGranted: voteGranted}, nil
}

func (handler *RequestHandler) HandleAppendEntriesRequest(request domain.AppendEntriesRequest) (*domain.AppendEntriesResponse, error) {
	// TODO: add checks related to log entries

	ctx := handler.ctx
	utils.CheckTerm(ctx, request.Term)

	log.Printf("Node \"%s\" received request of append-entries: %v", ctx.GetNodeId(), request)

	currentTerm := ctx.GetCurrentTerm()
	success := request.Term >= currentTerm

	if success {
		// Stable phase started
		ctx.ResetVotedFor()
	}

	return &domain.AppendEntriesResponse{Term: currentTerm, Success: success}, nil
}
