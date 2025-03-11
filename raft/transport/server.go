package transport

import (
	"distributed-algorithms/raft/domain"
)

type HandleRequestForVoteRequest func(request *domain.RequestVoteRequest) (*domain.RequestVoteResponse, error)

type HandleAppendEntriesRequest func(request *domain.AppendEntriesRequest) (*domain.AppendEntriesResponse, error)

type Server interface {
	Listen() error

	Shutdown() error
}

type ServerFactory interface {
	NewServer(
		address string,
		handleRequestForVoteRequest HandleRequestForVoteRequest,
		handleAppendEntriesRequest HandleAppendEntriesRequest,
	) (Server, error)
}
