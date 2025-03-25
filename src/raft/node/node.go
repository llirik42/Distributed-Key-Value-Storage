package node

import (
	"distributed-algorithms/src/config"
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/raft"
	"distributed-algorithms/src/raft/loops"
	"distributed-algorithms/src/raft/transport"
	"encoding/json"
	"fmt"
	logging "log"
)

func StartRaftNode(
	config config.RaftConfig,
	ctx *context.Context,
	raftServerFactory transport.ServerFactory,
	raftClientFactory transport.ClientFactory,
) error {
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
			logging.Printf("failed to shutdown RAFT-server gracefully: %s", err)
		}
	}(server)

	// Create connections to other nodes
	var clients []transport.Client
	for index, nodeAddress := range config.OtherNodes {
		newClient, errClient := raftClientFactory.NewClient(
			index,
			nodeAddress,
			messageHandler.HandleRequestVoteResponse,
			messageHandler.HandleAppendEntriesResponse,
		)

		if errClient != nil {
			logging.Fatalf("failed to creation connection to node %s: %s", nodeAddress, errClient)
		}

		defer func(newClient transport.Client) {
			if err := newClient.Close(); err != nil {
				logging.Printf("failed to close node-connection gracefully: %s", err)
			}
		}(newClient)

		clients = append(clients, newClient)
	}

	ctx.SetClients(clients)
	ctx.StartTickers()
	ctx.BecomeFollower()

	go loops.LeaderLoop(ctx)
	go loops.FollowerCandidateLoop(ctx)

	a, _ := json.MarshalIndent(config, "", " ")
	logging.Printf("node is starting with configuration %s", a)

	if errListen := server.Listen(); errListen != nil {
		return fmt.Errorf("node cannot start listen RAFT-connections: %w", errListen)
	}

	return nil
}
