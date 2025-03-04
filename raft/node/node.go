package node

import (
	"distributed-algorithms/config"
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/transport"
	"log"
)

type Node struct {
	ctx *context.Context
}

func NewNode(cfg config.RaftConfig, serverFactory transport.ServerFactory, clientFactory transport.ClientFactory) (*Node, error) {
	ctx := context.NewContext(cfg)
	requestHandler := raft.NewRequestHandler(ctx)

	server, err := serverFactory.NewServer(cfg.SelfNode.Address, requestHandler.HandleRequestVoteRequest, requestHandler.HandleAppendEntriesRequest)
	server
	if err != nil {
		return nil, err
	}

	var clients []transport.Client
	for _, nodeAddress := range cfg.OtherNodes {
		newClient, clientCreationErr := clientFactory.NewClient(nodeAddress)

		if clientCreationErr != nil {
			return nil, clientCreationErr
		}

		append(clients)
	}

	// Init leader-ticker
	//broadcastTimeoutMs := 100 // TODO
	//duration := getDurationMs(broadcastTimeoutMs)

	log.Println("Starting server")

	return nil, nil
}
