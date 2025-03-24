package common

import "distributed-algorithms/src/log"

type SetKeyRequest struct {
	Key   string
	Value any
}

type SetKeyResponse struct {
	Code     string
	LeaderId string
}

type GetKeyRequest struct {
	Key string
}

type GetKeyResponse struct {
	Value    any
	Code     string
	LeaderId string
}

type DeleteKeyRequest struct {
	Key string
}

type DeleteKeyResponse struct {
	Code     string
	LeaderId string
}

type GetClusterInfoRequest struct{}

type GetClusterInfoResponse struct {
	Code     string
	LeaderId string
	Info     struct {
		CurrentTerm uint32
		CommitIndex uint64
		LastApplied uint64
		NextIndex   []uint64
		MatchIndex  []uint64
	}
}

type GetLogRequest struct{}

type GetLogResponse struct {
	Code     string
	LeaderId string
	entries  []log.Entry
}
