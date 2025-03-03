package domain

type RequestVoteRequest struct {
	Term         int32
	CandidateId  uint32
	LastLogIndex int32
	LastLogTerm  int32
}
type RequestVoteResponse struct {
	Term        int32
	VoteGranted bool
}

type AppendEntriesRequest struct {
	Term         int32
	LeaderId     uint32
	PrevLogIndex int32
	PrevLogTerm  int32
	LeaderCommit int32
}

type AppendEntriesResponse struct {
	Term    int32
	Success bool
}
