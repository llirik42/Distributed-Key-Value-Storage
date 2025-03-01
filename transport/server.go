package transport

import "distributed-algorithms/dto"

type HandleRequestForVoteRequest func(request dto.RequestVoteRequest) (*dto.RequestVoteResponse, error)
type HandleAppendEntriesRequest func(request dto.AppendEntriesRequest) (*dto.AppendEntriesResponse, error)

type Server interface {
	Listen() error
}

type ServerFactory interface {
	NewServer(address string, handleRequestForVoteRequest HandleRequestForVoteRequest, handleAppendEntriesRequest HandleAppendEntriesRequest) (*Server, error)
}
