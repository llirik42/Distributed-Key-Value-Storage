package raft

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/log"
	"distributed-algorithms/src/raft/dto"
	"distributed-algorithms/src/raft/transport"
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
	ctx.Lock()
	defer ctx.Unlock()

	logStorage := ctx.GetLogStorage()
	a, _ := json.MarshalIndent(request, "", " ")
	logging.Printf("Node \"%s\" received request of vote: %s\n", ctx.GetNodeId(), a)

	checkTerm(ctx, request.Term) // TODO: Check this in gRPC-interceptor

	// Compare request's term and current term
	currentTerm := ctx.GetCurrentTerm()
	if request.Term < currentTerm {
		return &dto.RequestVoteResponse{Term: currentTerm, VoteGranted: false}, nil
	}

	// Compare current node's log and candidate's log
	currentNodeLastLogEntryMedata := logStorage.GetLastEntryMetadata()
	isCurrentNodeLogMoreUpToDate := log.CompareEntries(
		currentNodeLastLogEntryMedata.Term,
		currentNodeLastLogEntryMedata.Index,
		request.LastLogTerm,
		request.LastLogIndex,
	) < 0
	if isCurrentNodeLogMoreUpToDate {
		return &dto.RequestVoteResponse{Term: currentTerm, VoteGranted: false}, nil
	}

	// Try to vote for the candidate
	voteGranted := ctx.Vote(request.CandidateId)

	if voteGranted {
		ctx.ResetElectionTimeout()
	}

	return &dto.RequestVoteResponse{Term: currentTerm, VoteGranted: voteGranted}, nil
}

func (handler *MessageHandler) HandleAppendEntriesRequest(
	request *dto.AppendEntriesRequest,
) (*dto.AppendEntriesResponse, error) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	logStorage := ctx.GetLogStorage()
	a, _ := json.MarshalIndent(request, "", " ")
	logging.Printf("Node \"%s\" received request of append-entries: %s\n", ctx.GetNodeId(), a)

	checkTerm(ctx, request.Term) // TODO: Check this in gRPC-interceptor

	// Compare request's term and current term
	currentTerm := ctx.GetCurrentTerm()
	if request.Term < currentTerm {
		return &dto.AppendEntriesResponse{Term: currentTerm, Success: false}, nil
	}

	// Compare leader's log and current node's log
	logEntryMetadata, exists := logStorage.TryGetEntryMetadata(request.PrevLogIndex)
	logEntryTerm := logEntryMetadata.Term
	if !exists || logEntryTerm != request.PrevLogTerm {
		return &dto.AppendEntriesResponse{Term: currentTerm, Success: false}, nil
	}

	// Add new entries + resolve conflicts
	for i, entry := range request.Entries {
		newEntryIndex := request.PrevLogIndex + uint64(i) + 1
		logStorage.AddLogEntry(entry, newEntryIndex)
	}

	// Update commitIndex
	if request.LeaderCommit > ctx.GetCommitIndex() {
		newCommitIndex := min(request.LeaderCommit, logStorage.GetLastEntryMetadata().Index)
		ctx.SetCommitIndex(newCommitIndex)
	}

	ctx.BecomeFollower()
	ctx.SetLeaderId(request.LeaderId)

	return &dto.AppendEntriesResponse{Term: currentTerm, Success: true}, nil
}

func (handler *MessageHandler) HandleRequestVoteResponse(
	_ transport.Client,
	response *dto.RequestVoteResponse,
) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	a, _ := json.MarshalIndent(response, "", " ")
	logging.Printf("Node \"%s\" received response of vote: %s\n", ctx.GetNodeId(), a)

	checkTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor

	if !response.VoteGranted {
		return
	}

	// Got new vote
	voteNumber := ctx.IncrementVoteNumber()
	clusterSize := ctx.GetClusterSize()

	if voteNumber > clusterSize/2 {
		ctx.BecomeLeader()
		utils.SendAppendEntries(ctx)
	}
}

func (handler *MessageHandler) HandleAppendEntriesResponse(
	client transport.Client,
	response *dto.AppendEntriesResponse,
) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	a, _ := json.MarshalIndent(response, "", " ")
	logging.Printf("Node \"%s\" received response of append-entries: %s\n", ctx.GetNodeId(), a)

	checkTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor

	clientIndex := client.GetIndex() // Index of responder's connection

	if response.Success {
		lastSentIndex := ctx.GetLastSentIndex(clientIndex)
		ctx.SetNextIndex(clientIndex, lastSentIndex+1)
		ctx.SetMatchIndex(clientIndex, lastSentIndex)
	} else {
		ctx.DecrementNextIndex(clientIndex)
	}
}

func checkTerm(ctx *context.Context, term uint32) {
	if term > ctx.GetCurrentTerm() {
		ctx.SetCurrentTerm(term)
		ctx.BecomeFollower()
	}
}
