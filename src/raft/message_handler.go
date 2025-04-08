package raft

import (
	"distributed-key-value-storage/src/context"
	"distributed-key-value-storage/src/log"
	"distributed-key-value-storage/src/raft/dto"
	"distributed-key-value-storage/src/raft/transport"
	"distributed-key-value-storage/src/raft/utils"
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
	jsonMessage, _ := json.MarshalIndent(request, "", " ")
	logging.Printf(
		"Node \"%s\" received request of vote: %s\n",
		ctx.GetNodeId(),
		jsonMessage,
	)

	checkTerm(ctx, request.Term) // TODO: Check this in gRPC-interceptor

	// Compare request's term and current term
	currentTerm := ctx.GetCurrentTerm()
	if request.Term < currentTerm {
		return &dto.RequestVoteResponse{
			Term:        currentTerm,
			VoteGranted: false,
		}, nil
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
		return &dto.RequestVoteResponse{
			Term:        currentTerm,
			VoteGranted: false,
		}, nil
	}

	// Try to vote for the candidate
	voteGranted := ctx.Vote(request.CandidateId)

	if voteGranted {
		ctx.ResetElectionTimeout()
	}

	return &dto.RequestVoteResponse{
		Term:        currentTerm,
		VoteGranted: voteGranted,
	}, nil
}

func (handler *MessageHandler) HandleAppendEntriesRequest(
	request *dto.AppendEntriesRequest,
) (*dto.AppendEntriesResponse, error) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	logStorage := ctx.GetLogStorage()

	// Log request
	if len(request.Entries) > 0 {
		message := struct {
			Term         uint32
			LeaderId     string
			PrevLogIndex uint64
			PrevLogTerm  uint32
			EntriesCount int
			LeaderCommit uint64
		}{
			Term:         request.Term,
			LeaderId:     request.LeaderId,
			PrevLogIndex: request.PrevLogIndex,
			PrevLogTerm:  request.PrevLogTerm,
			EntriesCount: len(request.Entries),
			LeaderCommit: request.LeaderCommit,
		}

		jsonMessage, _ := json.MarshalIndent(message, "", " ")
		logging.Printf(
			"Node \"%s\" received request of append-entries: %s\n",
			ctx.GetNodeId(),
			jsonMessage,
		)
	}

	checkTerm(ctx, request.Term) // TODO: Check this in gRPC-interceptor

	// Compare request's term and current term
	currentTerm := ctx.GetCurrentTerm()
	if request.Term < currentTerm {
		return &dto.AppendEntriesResponse{Term: currentTerm, Success: false}, nil
	}

	logEntryMetadata, exists := logStorage.TryGetEntryMetadata(request.PrevLogIndex)

	if !exists {
		// Follower doesn't have log entry at prevLogIndex
		return &dto.AppendEntriesResponse{
			Term:          currentTerm,
			Success:       false,
			ConflictTerm:  0,
			ConflictIndex: logStorage.GetLength() + 1,
		}, nil
	}

	logEntryTerm := logEntryMetadata.Term
	if logEntryTerm != request.PrevLogTerm {
		// Term at prevLogIndex doesn't match prevLogTerm

		index, _ := logStorage.FindFirstEntryWithTerm(logEntryTerm)

		return &dto.AppendEntriesResponse{
			Term:          currentTerm,
			Success:       false,
			ConflictTerm:  logEntryTerm,
			ConflictIndex: index,
		}, nil
	}

	// Add new entries + resolve conflicts
	for i, entry := range request.Entries {
		newEntryIndex := request.PrevLogIndex + uint64(i) + 1
		logStorage.AddLogEntry(entry, newEntryIndex)
	}

	// Update commitIndex
	if request.LeaderCommit > ctx.GetCommitIndex() {
		newCommitIndex := min(request.LeaderCommit, logStorage.GetLength())
		ctx.SetCommitIndex(newCommitIndex)
	}

	ctx.BecomeFollower()
	ctx.SetLeaderId(request.LeaderId)

	return &dto.AppendEntriesResponse{
		Term:    currentTerm,
		Success: true,
	}, nil
}

func (handler *MessageHandler) HandleRequestVoteResponse(
	client transport.Client,
	response *dto.RequestVoteResponse,
) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	jsonMessage, _ := json.MarshalIndent(response, "", " ")
	logging.Printf(
		"Node \"%s\" received response of vote from \"%s\": %s\n",
		ctx.GetNodeId(),
		client.GetAddress(),
		jsonMessage,
	)

	checkTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor

	if !response.VoteGranted {
		return
	}

	// Got new vote
	voteNumber := ctx.IncrementVoteNumber()
	clusterSize := ctx.GetClusterSize()

	if voteNumber > clusterSize/2 {
		ctx.BecomeLeader()
		logging.Printf("Became leader!\n")
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

	if !response.Success {
		jsonMessage, _ := json.MarshalIndent(response, "", " ")
		logging.Printf(
			"Node \"%s\" received response of append-entries: %s\n",
			ctx.GetNodeId(),
			jsonMessage,
		)
	}

	if !checkTerm(ctx, response.Term) {
		// Received response from node with greater term
		return
	}

	clientIndex := client.GetIndex() // Index of responder's connection

	if response.Success {
		lastSentIndex := ctx.GetLastSentIndex(clientIndex)
		ctx.SetNextIndex(clientIndex, lastSentIndex+1)
		ctx.SetMatchIndex(clientIndex, lastSentIndex)
	} else {
		if response.ConflictTerm == 0 {
			// Conflict of index, so received response from node that doesn't have log entry at prevLogIndex
			ctx.SetNextIndex(clientIndex, response.ConflictIndex)
		} else {
			// Conflict of term, so received response from node which log[prevLogIndex] doesn't match prevLogTerm

			index, ok := ctx.GetLogStorage().FindLastEntryWithTerm(response.ConflictTerm)

			// Whether leader contains log entries with term ConflictTerm
			if ok {
				ctx.SetNextIndex(clientIndex, index+1)
			} else {
				ctx.SetNextIndex(clientIndex, response.ConflictIndex)
			}
		}
	}
}

func checkTerm(ctx *context.Context, term uint32) bool {
	if term > ctx.GetCurrentTerm() {
		ctx.SetCurrentTerm(term)
		ctx.BecomeFollower()
		return false
	}

	return true
}
