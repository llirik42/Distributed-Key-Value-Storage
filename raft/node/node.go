package node

import (
	"distributed-algorithms/raft/transport"
	"log"
	"sync"
	"sync/atomic"
)

type Node struct {
	nodeId int32

	nodeType      int
	nodeTypeMutex sync.Mutex

	currentTerm      atomic.Int32
	currentTermMutex sync.Mutex

	server transport.Server

	clients []transport.Client
}

func NewNode(serverFactory transport.ServerFactory, factory transport.ClientFactory) (*Node, error) {
	_, err := serverFactory.NewServer("0.0.0.0:8080", nil, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Starting server")

	return nil, nil
}

func (node *Node) Start() error {
	return node.server.Listen()
}

func (node *Node) GetId() int32 {
	return node.nodeId
}

func (node *Node) SetCurrentTerm(value int32) {
	node.currentTerm.Store(value)
}

func (node *Node) IncrementCurrentTerm() {
	node.currentTerm.Add(1)
}

func (node *Node) GetCurrentTerm() int32 {
	return node.currentTerm.Load()
}

func (node *Node) GetClients() []transport.Client {
	return node.clients
}

func (node *Node) BecomeFollower() {
	node.setType(follower)
}

func (node *Node) BecomeCandidate() {
	node.setType(candidate)
}

func (node *Node) BecomeLeader() {
	node.setType(leader)
}

func (node *Node) IsFollower() bool {
	return node.hasType(follower)
}

func (node *Node) IsCandidate() bool {
	return node.hasType(candidate)
}

func (node *Node) IsLeader() bool {
	return node.hasType(leader)
}

func (node *Node) setType(target int) {
	node.nodeTypeMutex.Lock()
	node.nodeType = target
	node.nodeTypeMutex.Unlock()
}

func (node *Node) hasType(target int) bool {
	node.nodeTypeMutex.Lock()
	result := node.nodeType == target
	node.nodeTypeMutex.Unlock()
	return result
}
