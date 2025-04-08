package main

import (
	"distributed-key-value-storage/src/client-interaction"
	"distributed-key-value-storage/src/config"
	"distributed-key-value-storage/src/context"
	kv "distributed-key-value-storage/src/key-value/in-memory"
	"distributed-key-value-storage/src/log/executor"
	log "distributed-key-value-storage/src/log/in-memory"
	"distributed-key-value-storage/src/raft/node"
	"distributed-key-value-storage/src/raft/transport/grpc"
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
	commandExecutor := executor.NewCommandExecutor(ctx, cfg.RaftConfig.ExecutedCommandsKey)
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
