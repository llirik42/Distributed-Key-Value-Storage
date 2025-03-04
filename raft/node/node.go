package node

import (
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/transport"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	currentTerm atomic.Int32

	voted         bool // TODO В какой момент нужно сбросить?
	votedFor      uint32
	votedForMutex sync.Mutex

	voteNumber atomic.Int32
	id         uint32

	role      int
	roleMutex sync.Mutex

	candidateLoopTicker *time.Ticker
	leaderLoopTicker    *time.Ticker

	server  transport.Server
	clients []transport.Client
}

func NewNode(serverFactory transport.ServerFactory, factory transport.ClientFactory) (*Node, error) {
	handler := raft.NewRequestHandler(ctx)

	_, err := serverFactory.NewServer("0.0.0.0:8080", handler.HandleRequestVoteRequest, nil)
	if err != nil {
		return nil, err
	}

	// Init leader-ticker
	//broadcastTimeoutMs := 100 // TODO
	//duration := getDurationMs(broadcastTimeoutMs)

	log.Println("Starting server")

	return nil, nil
}
