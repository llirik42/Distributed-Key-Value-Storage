package transport

import (
	"distributed-algorithms/src/raft/dto"
)

type HandleRequestForVoteResponse func(client Client, response *dto.RequestVoteResponse)

type HandleAppendEntriesResponse func(client Client, response *dto.AppendEntriesResponse)

// Client TODO: Split into 3 interfaces (SendRequestForVote + SendAppendEntries, Close, GetIndex?)
type Client interface {
	SendRequestForVote(request dto.RequestVoteRequest) error

	SendAppendEntries(request dto.AppendEntriesRequest) error

	GetIndex() int

	Close() error
}

type ClientFactory interface {
	NewClient(
		index int,
		address string,
		handleRequestForVoteResponse HandleRequestForVoteResponse,
		handleAppendEntriesResponse HandleAppendEntriesResponse,
	) (Client, error)
}
