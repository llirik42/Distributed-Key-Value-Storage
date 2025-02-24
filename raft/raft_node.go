package raft

import (
	"log"
	"math/rand"
	"time"
)

type Node struct {
	nodeState   int
	commonState CommonState
	leaderState LeaderState
	// add field for client connections
	isElectionTimerReset chan bool
}

func electionTimer(node *Node) {
	minElectionTimeoutMs := 500
	maxElectionTimeoutMs := 700

	electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs

	duration := time.Duration(electionTimeoutMs) * time.Millisecond

	timer := time.NewTimer(duration)
	log.Println("1", duration)
	<-timer.C
	log.Println("2", duration)
	timer.Stop()
	timer.Reset(duration)
	log.Println("Timer elapsed")
}

func NewNode() *Node {
	node := new(Node)
	node.nodeState = FOLLOWER
	node.commonState = CommonState{
		currentTerm: 0,
		votedFor:    noVotedFor,
		commitIndex: 0,
		lastApplied: 0,
	}

	return node
}

func (node *Node) resetElectionTimeout() error {
	node.isElectionTimerReset = true // Возможна гонка данных
	return nil
}
