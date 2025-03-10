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

	cfg, err1 := config.NewConfiguration(filePath)
	if err1 != nil {
		log.Fatal(err1)
	}

	err2 := node.StartRaftNode(cfg.RaftConfig, grpc.ServerFactory{}, grpc.ClientFactory{})
	if err2 != nil {
		log.Fatal(err2)
	}
}
