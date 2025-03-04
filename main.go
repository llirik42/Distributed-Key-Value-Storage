package main

import (
	"distributed-algorithms/config"
	"distributed-algorithms/raft/node"
	"distributed-algorithms/raft/transport/grpc"
	"log"
)

func main() {
	cfg1 := config.Config{RaftConfig: config.RaftConfig{
		BroadcastTimeMs:      1000,
		MinElectionTimeoutMs: 2000,
		MaxElectionTimeoutMs: 3000,
		OtherNodes:           []string{"localhost:8002"},
		SelfNode: struct {
			Address string `env:"SELF_ADDRESS,required"`
			Id      string `env:"SELF_ID,required"`
		}{Address: "localhost:8001", Id: "Node-1"},
	}}
	go func() {
		err := node.StartNode(cfg1, grpc.ServerFactory{}, grpc.ClientFactory{})
		if err != nil {
			log.Fatal(err)
		}
	}()

	cfg2 := config.Config{RaftConfig: config.RaftConfig{
		BroadcastTimeMs:      1000,
		MinElectionTimeoutMs: 2000,
		MaxElectionTimeoutMs: 3000,
		OtherNodes:           []string{"localhost:8001"},
		SelfNode: struct {
			Address string `env:"SELF_ADDRESS,required"`
			Id      string `env:"SELF_ID,required"`
		}{Address: "localhost:8002", Id: "Node-2"},
	}}
	err := node.StartNode(cfg2, grpc.ServerFactory{}, grpc.ClientFactory{})
	if err != nil {
		log.Fatal("Node-2: ", err)
	}
}
