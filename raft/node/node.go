package node

import (
	"distributed-algorithms/config"
	"distributed-algorithms/context"
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/loops"
	"distributed-algorithms/raft/transport"
	"encoding/json"
	"log"
)

func StartNode(config config.Config, raftServerFactory transport.ServerFactory, raftClientFactory transport.ClientFactory) error {
	err := startRaftNode(config.RaftConfig, raftServerFactory, raftClientFactory)

	if err != nil {
		// TODO: handle
		return err
	}

	return nil
}

func startRaftNode(config config.RaftConfig, raftServerFactory transport.ServerFactory, raftClientFactory transport.ClientFactory) error {
	ctx := context.NewContext(config)
	messageHandler := raft.NewMessageHandler(ctx)

	server, err := raftServerFactory.NewServer(config.SelfNode.Address, messageHandler.HandleRequestVoteRequest, messageHandler.HandleAppendEntriesRequest)
	if err != nil {
		// TODO: handle error
		return err
	}
	defer func(server transport.Server) {
		err := server.Shutdown()
		if err != nil {
			// TODO: handle error
		}
	}(server)

	var clients []transport.Client

	// Create connections to other nodes
	for _, nodeAddress := range config.OtherNodes {
		newClient, clientCreationErr := raftClientFactory.NewClient(nodeAddress, messageHandler.HandleRequestVoteResponse, messageHandler.HandleAppendEntriesResponse)

		if clientCreationErr != nil {
			// TODO: handle error
		}

		defer func(newClient transport.Client) {
			err := newClient.Close()
			if err != nil {
				// TODO: handle error
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

	listenErr := server.Listen()
	if listenErr != nil {
		// TODO: handle error
		return listenErr
	}

	return nil
}
