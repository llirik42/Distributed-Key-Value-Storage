package domain

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	Healthy bool
}

type RequestVoteRequest struct {
	Term         int32
	CandidateId  string
	LastLogIndex int32
	LastLogTerm  int32
}
type RequestVoteResponse struct {
	Term        int32
	VoteGranted bool
}

type AppendEntriesRequest struct {
	Term         int32
	LeaderId     string
	PrevLogIndex int32
	PrevLogTerm  int32
	LeaderCommit int32
}

type AppendEntriesResponse struct {
	Term    int32
	Success bool
}
