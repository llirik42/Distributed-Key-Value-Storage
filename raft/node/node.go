package node

import (
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/transport"
	"log"
)

type Node struct {
}

func NewNode(serverFactory transport.ServerFactory, factory transport.ClientFactory) (*Node, error) {
	var node *Node = nil

	handler := raft.NewRequestHandler(node)

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
