package main

import (
	"distributed-algorithms/config"
	"distributed-algorithms/raft/node"
	"distributed-algorithms/raft/transport/grpc"
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
