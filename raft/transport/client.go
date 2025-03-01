package transport

import (
	"distributed-algorithms/raft/dto"
)

type HandleRequestForVoteResponse func(request dto.RequestVoteResponse) error
type HandleAppendEntriesResponse func(request dto.AppendEntriesResponse) error

type Client interface {
	SendRequestForVote(request dto.RequestVoteRequest) (*dto.RequestVoteResponse, error)

	SendAppendEntries(request dto.AppendEntriesRequest) (*dto.AppendEntriesResponse, error)

	Close() error
}

type ClientFactory interface {
	NewClient(address string) (*Client, error)
}
