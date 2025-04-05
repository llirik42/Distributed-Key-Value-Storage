package dto

import (
	"distributed-algorithms/src/log"
)

type RequestVoteRequest struct {
	Term         uint32
	CandidateId  string
	LastLogIndex uint64
	LastLogTerm  uint32
}
type RequestVoteResponse struct {
	Term        uint32
	VoteGranted bool
}

type AppendEntriesRequest struct {
	Term         uint32
	LeaderId     string
	PrevLogIndex uint64
	PrevLogTerm  uint32
	Entries      []log.Entry
	LeaderCommit uint64
}

type AppendEntriesResponse struct {
	Term          uint32
	Success       bool
	ConflictTerm  uint32
	ConflictIndex uint64
}
