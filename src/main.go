package main

import (
	"distributed-algorithms/src/config"
	"distributed-algorithms/src/raft/node"
	"distributed-algorithms/src/raft/transport/grpc"
	"log"
	"os"
)

func main() {
	args := os.Args
	filePath := args[1]

	cfg, errConfig := config.NewConfiguration(filePath)
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	if err := node.StartRaftNode(cfg.RaftConfig, grpc.ServerFactory{}, grpc.ClientFactory{}); err != nil {
		log.Fatal(err)
	}
}
