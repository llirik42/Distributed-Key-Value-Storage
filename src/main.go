package main

import (
	"distributed-algorithms/src/client-interaction/common"
	"distributed-algorithms/src/client-interaction/restapi"
	"distributed-algorithms/src/config"
	"distributed-algorithms/src/context"
	kv "distributed-algorithms/src/key-value/in-memory"
	log "distributed-algorithms/src/log/in-memory"
	"distributed-algorithms/src/raft/node"
	"distributed-algorithms/src/raft/transport/grpc"
	logging "log"
	"os"
)

func main() {
	args := os.Args
	filePath := args[1]

	cfg, err := config.NewConfiguration(filePath)
	if err != nil {
		logging.Fatalf("error loading configuration: %v", err)
	}

	ctx := context.NewContext(cfg.RaftConfig)
	ctx.SetKeyValueStorage(kv.NewStorage())
	ctx.SetLogStorage(log.NewStorage())

	serverFactory := grpc.NewServerFactory()
	clientFactory := grpc.NewClientFactory()

	go func() {
		if err := node.StartRaftNode(cfg.RaftConfig, ctx, serverFactory, clientFactory); err != nil {
			logging.Fatalf("error starting RAFT-node: %v", err)
		}
	}()

	requestHandler := common.NewRequestHandler(ctx)

	if err := restapi.StartServer(requestHandler, cfg.RestConfig); err != nil {
		logging.Fatalf("error starting restapi: %v", err)
	}
}
