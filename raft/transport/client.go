package transport

import (
	"distributed-algorithms/raft/domain"
)

type HandleRequestForVoteResponse func(response *domain.RequestVoteResponse)

type HandleAppendEntriesResponse func(response *domain.AppendEntriesResponse)

type Client interface {
	SendRequestForVote(request domain.RequestVoteRequest) error

	SendAppendEntries(request domain.AppendEntriesRequest) error

	SendHealthCheck(request domain.HealthCheckRequest) (*domain.HealthCheckResponse, error)

	Close() error
}

type ClientFactory interface {
	NewClient(
		address string,
		handleRequestForVoteResponse HandleRequestForVoteResponse,
		handleAppendEntriesResponse HandleAppendEntriesResponse) (Client, error)
}
