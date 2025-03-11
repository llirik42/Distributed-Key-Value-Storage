package node

import (
	"distributed-algorithms/config"
	"distributed-algorithms/context"
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/loops"
	"distributed-algorithms/raft/transport"
	"encoding/json"
	"fmt"
	"log"
)

func StartRaftNode(
	config config.RaftConfig,
	raftServerFactory transport.ServerFactory,
	raftClientFactory transport.ClientFactory,
) error {
	ctx := context.NewContext(config)
	messageHandler := raft.NewMessageHandler(ctx)

	// Create and start server
	server, errServer := raftServerFactory.NewServer(
		config.SelfNode.Address,
		messageHandler.HandleRequestVoteRequest,
		messageHandler.HandleAppendEntriesRequest,
	)
	if errServer != nil {
		return fmt.Errorf("failed to RAFT-server for node: %w", errServer)
	}
	defer func(server transport.Server) {
		if err := server.Shutdown(); err != nil {
			log.Printf("failed to shutdown RAFT-server gracefully: %s", err)
		}
	}(server)

	// Create connections to other nodes
	var clients []transport.Client
	for _, nodeAddress := range config.OtherNodes {
		newClient, errClient := raftClientFactory.NewClient(
			nodeAddress,
			messageHandler.HandleRequestVoteResponse,
			messageHandler.HandleAppendEntriesResponse,
		)

		if errClient != nil {
			log.Fatalf("failed to creation connection to node %s: %s", nodeAddress, errClient)
		}

		defer func(newClient transport.Client) {
			if err := newClient.Close(); err != nil {
				log.Printf("failed to close node-connection gracefully: %s", err)
			}
		}(newClient)

		clients = append(clients, newClient)
	}

	ctx.SetServer(&server)
	ctx.SetClients(clients)

	ctx.StartTickers()
	ctx.BecomeFollower()
	go loops.LeaderLoop(ctx)
	go loops.FollowerCandidateLoop(ctx)

	a, _ := json.MarshalIndent(config, "", " ")
	log.Printf("Node is starting with configuration %s\n", a)

	if errListen := server.Listen(); errListen != nil {
		return fmt.Errorf("node cannot start listen RAFT-connections: %w", errListen)
	}

	return nil
}
