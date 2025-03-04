package transport

import (
	"distributed-algorithms/raft/domain"
)

type HandleRequestForVoteResponse func(request domain.RequestVoteResponse) error
type HandleAppendEntriesResponse func(request domain.AppendEntriesResponse) error

type Client interface {
	SendRequestForVote(request domain.RequestVoteRequest) (*domain.RequestVoteResponse, error)

	SendAppendEntries(request domain.AppendEntriesRequest) (*domain.AppendEntriesResponse, error)

	Close() error
}

type ClientFactory interface {
	NewClient(address string) (Client, error)
}
