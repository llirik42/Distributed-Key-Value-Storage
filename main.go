package main

import (
	"distributed-algorithms/config"
	"distributed-algorithms/raft/node"
	"distributed-algorithms/raft/transport/grpc"
	"log"
)

func main() {
	cfg1, err11 := config.NewConfiguration(".env")
	if err11 != nil {
		// TODO: handle
		log.Fatal(err11)
	}

	err12 := node.StartNode(*cfg1, grpc.ServerFactory{}, grpc.ClientFactory{})
	if err12 != nil {
		log.Fatal(err12)
	}
	
}
