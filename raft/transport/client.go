package transport

import (
	"distributed-algorithms/raft/domain"
)

type HandleRequestForVoteResponse func(response *domain.RequestVoteResponse) error
type HandleAppendEntriesResponse func(response *domain.AppendEntriesResponse) error

type Client interface {
	SendRequestForVote(request domain.RequestVoteRequest) error

	SendAppendEntries(request domain.AppendEntriesRequest) error

	Close() error
}

type ClientFactory interface {
	NewClient(address string, handleRequestForVoteResponse HandleRequestForVoteResponse, handleAppendEntriesResponse HandleAppendEntriesResponse) (Client, error)
}
