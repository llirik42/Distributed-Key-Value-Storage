package main

import (
	"distributed-algorithms/raft/node"
	"distributed-algorithms/raft/transport/grpc"
	"fmt"
)

func main() {
	serverFactory := grpc.ServerFactory{}
	clientFactory := grpc.ClientFactory{}

	_, err := node.NewNode(serverFactory, clientFactory)

	if err != nil {
		fmt.Println(err)
	}
}
