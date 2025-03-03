package node

import (
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/transport"
	"log"
)

type Node struct {
	ctx context.Context
}

func NewNode(serverFactory transport.ServerFactory, factory transport.ClientFactory) (*Node, error) {
	ctx := context.NewContext()

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
