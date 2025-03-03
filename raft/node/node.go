package node

import (
	"distributed-algorithms/raft/transport"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type HeartbeatInfo struct {
	time     time.Time
	leaderId Id
}

type Id int32

type Node struct {
	nodeId   Id
	votedFor Id

	voteNumber atomic.Int32

	nodeType      int
	nodeTypeMutex sync.Mutex

	currentTerm atomic.Int32

	candidateLoopTicker *time.Ticker
	leaderLoopTicker    *time.Ticker

	server transport.Server

	clients []transport.Client
}

func NewNode(serverFactory transport.ServerFactory, factory transport.ClientFactory) (*Node, error) {
	_, err := serverFactory.NewServer("0.0.0.0:8080", nil, nil)
	if err != nil {
		return nil, err
	}

	// Init leader-ticker
	//broadcastTimeoutMs := 100 // TODO
	//duration := getDurationMs(broadcastTimeoutMs)

	log.Println("Starting server")

	return nil, nil
}

func (node *Node) Start() error {
	node.nodeType = follower // is it necessary to use node.setType?

	return node.server.Listen()
}

func (node *Node) IncrementVoteNumber() {
	node.voteNumber.Add(1)
}

func (node *Node) GetVoteNumber() int32 {
	return node.voteNumber.Load()
}

func (node *Node) ResetVoteNumber() {
	node.voteNumber.Store(1) // node votes for itself
}

func (node *Node) CheckTerm(term int32) {
	if term > node.GetCurrentTerm() {
		node.SetCurrentTerm(term)
		node.BecomeFollower()
	}
}

func (node *Node) GetId() Id {
	return node.nodeId
}

func (node *Node) SetCurrentTerm(value int32) {
	node.currentTerm.Store(value)
}

func (node *Node) IncrementCurrentTerm() int32 {
	return node.currentTerm.Add(1)
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
