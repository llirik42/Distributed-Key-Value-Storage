package in_memory

import "distributed-algorithms/src/log"

type Log struct {
	entries []log.Entry
}

func NewLog() *Log {
	return &Log{}
}
