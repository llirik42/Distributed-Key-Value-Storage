package node

import (
	"distributed-algorithms/config"
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/transport"
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
	requestHandler := raft.NewRequestHandler(ctx)

	server, err := raftServerFactory.NewServer(config.SelfNode.Address, requestHandler.HandleRequestVoteRequest, requestHandler.HandleAppendEntriesRequest)
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
		newClient, clientCreationErr := raftClientFactory.NewClient(nodeAddress)

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

	listenErr := server.Listen()
	if listenErr != nil {
		// TODO: handle error
		return listenErr
	}

	return nil
}
