package raft

import (
	"distributed-algorithms/raft/transport"
	"log"
)

type Node struct {
	Server transport.Server
}

func NewNode(serverFactory transport.ServerFactory) (*Node, error) {
	_, err := serverFactory.NewServer("0.0.0.0:8080", nil, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Starting server")

	return nil, nil
}

func (node *Node) Start() error {
	return node.Server.Listen()
}
