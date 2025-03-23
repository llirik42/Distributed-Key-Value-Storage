package raft

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/log"
	"distributed-algorithms/src/raft/dto"
	"distributed-algorithms/src/raft/utils"
	"encoding/json"
	logging "log"
)

type MessageHandler struct {
	ctx *context.Context
}

func NewMessageHandler(ctx *context.Context) *MessageHandler {
	return &MessageHandler{ctx: ctx}
}

func (handler *MessageHandler) HandleRequestVoteRequest(
	request *dto.RequestVoteRequest,
) (*dto.RequestVoteResponse, error) {
	ctx := handler.ctx

	a, _ := json.MarshalIndent(request, "", " ")
	logging.Printf("Node \"%s\" received request of vote: %s\n", ctx.GetNodeId(), a)

	ctx.Lock()
	defer ctx.Unlock()

	checkTerm(ctx, request.Term) // TODO: Check this in gRPC-interceptor
	currentTerm := ctx.GetCurrentTerm()

	if request.Term < currentTerm {
		return &dto.RequestVoteResponse{Term: currentTerm, VoteGranted: false}, nil
	}

	// Compare current node's log and candidate's log
	currentNodeLastLogEntryMedata := ctx.GetLastLogEntryMetadata()
	isCurrentNodeLogMoreUpToDate := log.CompareEntries(
		currentNodeLastLogEntryMedata.Term,
		currentNodeLastLogEntryMedata.Index,
		request.LastLogTerm,
		request.LastLogIndex,
	) < 0
	if isCurrentNodeLogMoreUpToDate {
		return &dto.RequestVoteResponse{Term: currentTerm, VoteGranted: false}, nil
	}

	voteGranted := ctx.Vote(request.CandidateId)

	if voteGranted {
		ctx.ResetNewElectionTimeout()
	}

	return &dto.RequestVoteResponse{Term: currentTerm, VoteGranted: voteGranted}, nil
}

func (handler *MessageHandler) HandleAppendEntriesRequest(
	request *dto.AppendEntriesRequest,
) (*dto.AppendEntriesResponse, error) {
	// TODO: add checks related to log entries

	ctx := handler.ctx

	a, _ := json.MarshalIndent(request, "", " ")
	logging.Printf("Node \"%s\" received request of append-entries: %s\n", ctx.GetNodeId(), a)

	ctx.Lock()
	defer ctx.Unlock()

	checkTerm(ctx, request.Term) // TODO: Check this in gRPC-interceptor
	currentTerm := ctx.GetCurrentTerm()

	if request.Term < currentTerm {
		return &dto.AppendEntriesResponse{Term: currentTerm, Success: false}, nil
	}

	logEntryTerm, exists := ctx.GetLogEntryTerm(request.PrevLogIndex)
	if !exists || logEntryTerm != request.PrevLogTerm {
		return &dto.AppendEntriesResponse{Term: currentTerm, Success: false}, nil
	}

	var success bool
	if request.Term < currentTerm {
		success = false
	} else {

	}

	requestTerm := request.Term
	success := requestTerm >= currentTerm

	// "checkTerm" is implicitly called here
	if success {
		// Stable phase started
		ctx.BecomeFollower()
\	}

	return &dto.AppendEntriesResponse{Term: currentTerm, Success: success}, nil
}

func (handler *MessageHandler) HandleRequestVoteResponse(
	response *dto.RequestVoteResponse,
) {
	ctx := handler.ctx

	a, _ := json.MarshalIndent(response, "", " ")
	logging.Printf("Node \"%s\" received response of vote: %s\n", ctx.GetNodeId(), a)

	ctx.Lock()
	defer ctx.Unlock()

	checkTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor

	if !response.VoteGranted {
		return
	}

	// Got new vote
	voteNumber := ctx.IncrementVoteNumber()
	clusterSize := ctx.GetClusterSize()

	if voteNumber > clusterSize/2 {
		ctx.BecomeLeader()
		utils.SendHeartbeat(ctx)
	}
}

func (handler *MessageHandler) HandleAppendEntriesResponse(
	response *dto.AppendEntriesResponse,
) {
	// TODO: add checks related to logs

	ctx := handler.ctx

	a, _ := json.MarshalIndent(response, "", " ")
	logging.Printf("Node \"%s\" received response of append-entries: %s\n", ctx.GetNodeId(), a)

	ctx.Lock()
	defer ctx.Unlock()

	checkTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor
}

func checkTerm(ctx *context.Context, term uint32) {
	if term > ctx.GetCurrentTerm() {
		ctx.SetCurrentTerm(term)
		ctx.BecomeFollower()
	}
}
