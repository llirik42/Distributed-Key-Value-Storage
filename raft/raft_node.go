package raft

import (
	"distributed-algorithms/raft/test"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	nodeState      int
	nodeStateMutex sync.Mutex

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
	node.nodeState = test.FOLLOWER
	node.nodeStateMutex = sync.Mutex{}

	node.isElectionTimerReset = make(chan bool)

	node.commonState = CommonState{
		currentTerm: 0,
		votedFor:    noVotedFor,
		commitIndex: 0,
		lastApplied: 0,
	}
	node.voteMutex = sync.Mutex{}

	return node
}

func (node *Node) changeState(newState int) {
	node.voteMutex.Lock()
	node.nodeState = newState
	node.voteMutex.Unlock()
}

func (node *Node) BecomeLeader() {
	node.changeState(test.LEADER)
}

func (node *Node) BecomeCandidate() {
	node.changeState(test.CANDIDATE)
}

func (node *Node) BecomeFollower() {
	node.changeState(test.FOLLOWER)
}

func (node *Node) IsLeader() bool {
	node.voteMutex.Lock()
	result := node.nodeState == test.LEADER
	node.voteMutex.Unlock()
	return result
}

func (node *Node) IsCandidate() bool {
	node.voteMutex.Lock()
	result := node.nodeState == test.LEADER
	node.voteMutex.Unlock()
	return result
}

// IsFollower TODO: Is this method used anywhere?
func (node *Node) IsFollower() bool {
	node.voteMutex.Lock()
	result := node.nodeState == test.LEADER
	node.voteMutex.Unlock()
	return result
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
