package context

import (
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/transport"
	"sync"
	"sync/atomic"
	"time"
)

type Context struct {
	currentTerm atomic.Int32

	voted         bool // TODO В какой момент нужно сбросить?
	votedFor      uint32
	votedForMutex sync.Mutex
	voteNumber    atomic.Int32

	nodeId        uint32
	nodeRole      int
	nodeRoleMutex sync.Mutex

	candidateLoopTicker *time.Ticker
	leaderLoopTicker    *time.Ticker

	server  transport.Server
	clients []transport.Client
}

func NewContext() *Context {
	ctx := &Context{
		currentTerm:   atomic.Int32{},
		voted:         false,
		votedFor:      0, // Doesn't matter when voted = false
		votedForMutex: sync.Mutex{},
		voteNumber:    atomic.Int32{},
	}

	return ctx
}

func (ctx *Context) GetClients() []transport.Client {
	return ctx.clients
}

func (ctx *Context) GetNodeId() uint32 {
	return ctx.nodeId
}

func (ctx *Context) IsFollower() bool {
	return ctx.hasType(domain.FOLLOWER)
}

func (ctx *Context) IsCandidate() bool {
	return ctx.hasType(domain.CANDIDATE)
}

func (ctx *Context) IsLeader() bool {
	return ctx.hasType(domain.LEADER)
}

func (ctx *Context) setRole(target int) {
	ctx.nodeRoleMutex.Lock()
	ctx.nodeRole = target
	ctx.nodeRoleMutex.Unlock()
}

func (ctx *Context) hasType(target int) bool {
	ctx.nodeRoleMutex.Lock()
	result := ctx.nodeRole == target
	ctx.nodeRoleMutex.Unlock()
	return result
}

func (ctx *Context) ResetVote() {
	ctx.votedForMutex.Lock()
	ctx.voted = false
	ctx.votedForMutex.Unlock()
}

func (ctx *Context) Vote(candidateId uint32) bool {
	ctx.votedForMutex.Lock()

	var result bool

	if !ctx.voted {
		ctx.votedFor = candidateId
		ctx.voted = true
		result = true
	} else {
		// If the candidate we voted for falls, we'll again vote for it
		// If result = false, in this case we won't vote for the candidate
		result = ctx.votedFor == candidateId // In case ca
	}

	ctx.votedForMutex.Unlock()

	return result
}

func (ctx *Context) VoteForMyself() bool {
	result := ctx.Vote(ctx.GetNodeId())

}

func (ctx *Context) IncrementVoteNumber() {
	ctx.voteNumber.Add(1)
}

func (ctx *Context) GetVoteNumber() int32 {
	return ctx.voteNumber.Load()
}

func (ctx *Context) ResetVoteNumber() {
	ctx.voteNumber.Store(1) // node votes for itself
}

func (ctx *Context) CheckTerm(term int32) {
	if term > ctx.GetCurrentTerm() {
		ctx.SetCurrentTerm(term)
		ctx.BecomeFollower()
	}
}

func (ctx *Context) SetCurrentTerm(value int32) {
	ctx.currentTerm.Store(value)
}

func (ctx *Context) IncrementCurrentTerm() int32 {
	return ctx.currentTerm.Add(1)
}

func (ctx *Context) GetCurrentTerm() int32 {
	return ctx.currentTerm.Load()
}

func (ctx *Context) BecomeFollower() {
	ctx.setRole(domain.FOLLOWER)
	// TODO: do something with tickers
}

func (ctx *Context) BecomeCandidate() {
	ctx.setRole(domain.CANDIDATE)
	// TODO: do something with tickers
}

func (ctx *Context) BecomeLeader() {
	ctx.setRole(domain.LEADER)
	// TODO: do something with tickers}
}
