package node

import (
	"distributed-algorithms/config"
	"distributed-algorithms/context"
	"distributed-algorithms/raft"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/loops"
	"distributed-algorithms/raft/transport"
	"encoding/json"
	"log"
	"time"
)

func StartRaftNode(config config.RaftConfig, raftServerFactory transport.ServerFactory, raftClientFactory transport.ClientFactory) error {
	a, _ := json.MarshalIndent(config, "", " ")
	log.Printf("Node is starting with configuration:\n%s\n", a)

	ctx := context.NewContext(config)
	messageHandler := raft.NewMessageHandler(ctx)

	// Create and start server
	server, err := raftServerFactory.NewServer(config.SelfNode.Address, messageHandler.HandleRequestVoteRequest, messageHandler.HandleAppendEntriesRequest)
	if err != nil {
		// TODO: handle error
		return err
	}
	go serve(server)

	if ensuringErr := ensureServerStarted(config, raftClientFactory); ensuringErr != nil {
		// TODO: handle
	}
	log.Printf("Node is ready to accept connections!")

	// Create connections to other nodes
	var clients []transport.Client
	for _, nodeAddress := range config.OtherNodes {
		newClient, clientCreationErr := raftClientFactory.NewClient(
			nodeAddress,
			messageHandler.HandleRequestVoteResponse,
			messageHandler.HandleAppendEntriesResponse)

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
	loops.FollowerCandidateLoop(ctx)

	return nil
}

func ensureServerStarted(config config.RaftConfig, raftClientFactory transport.ClientFactory) error {
	// This client will only be used to send health-check,
	// so handleRequestForVoteResponse and handleAppendEntriesResponse can be nil

	client, clientCreationErr := raftClientFactory.NewClient(
		config.SelfNode.Address,
		nil,
		nil)

	if clientCreationErr != nil {
		return clientCreationErr
	}

	duration := time.Duration(config.HealthCheckTimeoutMs) * time.Millisecond
	ticker := time.NewTicker(duration)

	request := domain.HealthCheckRequest{}
	for range ticker.C {
		response, err := client.SendHealthCheck(request)

		if err != nil {
			log.Printf("Failed to send health check request: %v", err)
			// TODO: handle error
			continue
		}

		if response.Healthy {
			ticker.Stop()
			break
		}

		// TODO: handle unhealthy
	}

	return nil
}

func serve(server transport.Server) {
	defer func(server transport.Server) {
		err := server.Shutdown()
		if err != nil {
			// TODO: handle error
		}
	}(server)

	err := server.Listen()
	if err != nil {
		// TODO: handle
	}
}
