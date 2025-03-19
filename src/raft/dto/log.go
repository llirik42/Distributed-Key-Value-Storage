package dto

import "distributed-algorithms/src/domain"

type LogEntry struct {
	Index   uint64
	Term    uint32
	Command domain.Command
}
