package transport

import (
	"distributed-algorithms/src/raft/dto"
)

type HandleRequestForVoteResponse func(response *dto.RequestVoteResponse)

type HandleAppendEntriesResponse func(response *dto.AppendEntriesResponse)

type Client interface {
	SendRequestForVote(request dto.RequestVoteRequest) error

	SendAppendEntries(request dto.AppendEntriesRequest) error

	Close() error
}

type ClientFactory interface {
	NewClient(
		address string,
		handleRequestForVoteResponse HandleRequestForVoteResponse,
		handleAppendEntriesResponse HandleAppendEntriesResponse,
	) (Client, error)
}
