package main

import (
	"distributed-algorithms/src/client-interaction"
	"distributed-algorithms/src/config"
	"distributed-algorithms/src/context"
	kv "distributed-algorithms/src/key-value/in-memory"
	"distributed-algorithms/src/log/executor"
	log "distributed-algorithms/src/log/in-memory"
	"distributed-algorithms/src/raft/node"
	"distributed-algorithms/src/raft/transport/grpc"
	logging "log"
	"os"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("Invalid number of arguments")
	}

	filePath := args[1]

	cfg, err := config.NewConfiguration(filePath)
	if err != nil {
		logging.Fatalf("error loading configuration: %v", err)
	}

	ctx := context.NewContext(cfg.RaftConfig)
	keyValueStorage := kv.NewStorage()
	commandExecutor := executor.NewCommandExecutor(ctx)
	logStorage := log.NewStorage()

	ctx.SetKeyValueStorage(keyValueStorage)
	ctx.SetLogStorage(logStorage)
	ctx.SetCommandExecutor(commandExecutor)

	serverFactory := grpc.NewServerFactory()
	clientFactory := grpc.NewClientFactory()

	go func() {
		if err := node.StartRaftNode(cfg.RaftConfig, ctx, serverFactory, clientFactory); err != nil {
			logging.Fatalf("error starting RAFT-node: %v", err)
		}
	}()

	if err := client_interaction.StartServer(ctx, cfg.RestConfig); err != nil {
		logging.Fatalf("error starting restapi: %v", err)
	}
}
