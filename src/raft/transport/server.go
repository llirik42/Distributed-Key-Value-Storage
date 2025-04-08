package transport

import (
	"distributed-key-value-storage/src/raft/dto"
)

type HandleRequestForVoteRequest func(request *dto.RequestVoteRequest) (*dto.RequestVoteResponse, error)

type HandleAppendEntriesRequest func(request *dto.AppendEntriesRequest) (*dto.AppendEntriesResponse, error)

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
