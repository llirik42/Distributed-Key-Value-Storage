package raft

import (
	"distributed-algorithms/context"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/utils"
	"encoding/json"
	"log"
)

type MessageHandler struct {
	ctx *context.Context
}

func NewMessageHandler(ctx *context.Context) *MessageHandler {
	return &MessageHandler{ctx: ctx}
}

func (handler *MessageHandler) HandleRequestVoteRequest(request *domain.RequestVoteRequest) (*domain.RequestVoteResponse, error) {
	// TODO: add checks about candidate's log

	ctx := handler.ctx

	a, _ := json.MarshalIndent(request, "", " ")
	log.Printf("Node \"%s\" received request of vote: %s\n", ctx.GetNodeId(), a)

	checkTerm(ctx, request.Term)

	currentTerm := ctx.GetCurrentTerm()
	var voteGranted bool

	if request.Term < currentTerm {
		voteGranted = false
	} else {
		voteGranted = ctx.Vote(request.CandidateId)
	}

	if voteGranted {
		ctx.ResetNewElectionTimeout()
	}

	return &domain.RequestVoteResponse{Term: currentTerm, VoteGranted: voteGranted}, nil
}

func (handler *MessageHandler) HandleAppendEntriesRequest(request *domain.AppendEntriesRequest) (*domain.AppendEntriesResponse, error) {
	// TODO: add checks related to log entries

	ctx := handler.ctx

	a, _ := json.MarshalIndent(request, "", " ")
	log.Printf("Node \"%s\" received request of append-entries: %s\n", ctx.GetNodeId(), a)

	currentTerm := ctx.GetCurrentTerm()
	requestTerm := request.Term
	success := requestTerm >= currentTerm

	if success {
		// Stable phase started
		ctx.BecomeFollower()
		ctx.SetCurrentTerm(requestTerm)
	}

	return &domain.AppendEntriesResponse{Term: currentTerm, Success: success}, nil
}

func (handler *MessageHandler) HandleRequestVoteResponse(response *domain.RequestVoteResponse) error {
	ctx := handler.ctx

	a, _ := json.MarshalIndent(response, "", " ")
	log.Printf("Node \"%s\" received response of vote: %s\n", ctx.GetNodeId(), a)

	checkTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor

	if !response.VoteGranted {
		return nil
	}

	// Got new vote
	voteNumber := ctx.IncrementVoteNumber()
	clusterSize := ctx.GetClusterSize()

	if voteNumber > clusterSize/2 {
		ctx.BecomeLeader()
		utils.SendHeartbeat(ctx)
	}

	return nil
}

func (handler *MessageHandler) HandleAppendEntriesResponse(response *domain.AppendEntriesResponse) error {
	// TODO: add checks related to logs

	ctx := handler.ctx

	a, _ := json.MarshalIndent(response, "", " ")
	log.Printf("Node \"%s\" received response of append-entries: %s\n", ctx.GetNodeId(), a)

	checkTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor

	return nil
}

func checkTerm(ctx *context.Context, term int32) {
	if term > ctx.GetCurrentTerm() {
		ctx.SetCurrentTerm(term)
		ctx.BecomeFollower()
	}
}
