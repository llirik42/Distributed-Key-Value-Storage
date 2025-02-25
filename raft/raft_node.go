package raft

import (
	"distributed-algorithms/raft/test"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	nodeState   int
	commonState CommonState
	leaderState LeaderState
	// add field for client connections
	isElectionTimerReset chan bool
	voteMutex            sync.Mutex
}

func electionTimer(node *Node) {
	minElectionTimeoutMs := 500
	maxElectionTimeoutMs := 700

	electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs

	duration := time.Duration(electionTimeoutMs) * time.Millisecond

	time.Sleep(duration)

	isElectionTimerReset := <-node.isElectionTimerReset

	if isElectionTimerReset {
		node.isElectionTimerReset <- false
		electionTimer(node)
	} else {
		node.nodeState = test.CANDIDATE
		node.commonState.currentTerm++

		// Отправить всем RequestForVote
	}
}

func NewNode() *Node {
	node := new(Node)
	node.isElectionTimerReset = make(chan bool)
	node.nodeState = test.FOLLOWER
	node.commonState = CommonState{
		currentTerm: 0,
		votedFor:    noVotedFor,
		commitIndex: 0,
		lastApplied: 0,
	}
	node.voteMutex = sync.Mutex{}

	return node
}

func (node *Node) resetElectionTimeout() error {
	node.isElectionTimerReset <- true
	return nil
}

func (node *Node) GetCurrentTerm() int32 {
	return node.commonState.currentTerm // TODO: использовать тут канал, так как возможна гонка данных
}

func (node *Node) Vote(candidateId int32) bool {
	node.voteMutex.Lock()

	var result bool

	if node.commonState.votedFor == noVotedFor {
		node.commonState.votedFor = candidateId
		result = true
	} else {
		result = false
	}

	node.voteMutex.Unlock()

	return result
}
