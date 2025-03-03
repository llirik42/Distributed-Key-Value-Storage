package context

import (
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/transport"
	"sync"
	"sync/atomic"
	"time"
)

type Context struct {
	nodeId uint32

	voted         bool // TODO В какой момент нужно сбросить?
	votedFor      uint32
	votedForMutex sync.Mutex

	voteNumber atomic.Int32

	nodeType      int
	nodeTypeMutex sync.Mutex

	currentTerm atomic.Int32

	candidateLoopTicker *time.Ticker
	leaderLoopTicker    *time.Ticker

	clients []transport.Client
}

func NewContext() *Context {
	// TODO: init all fields

	ctx := &Context{}

	ctx.nodeType = domain.FOLLOWER // is it necessary to use ctx.setType here?
	return ctx
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

func (ctx *Context) IncrementVoteNumber() {
	ctx.voteNumber.Add(1)
}

func (ctx *Context) GetVoteNumber() int32 {
	return ctx.voteNumber.Load()
}

func (ctx *Context) ResetVoteNumber() {
	ctx.voteNumber.Store(1) // ctx votes for itself
}

func (ctx *Context) CheckTerm(term int32) {
	if term > ctx.GetCurrentTerm() {
		ctx.SetCurrentTerm(term)
		ctx.BecomeFollower()
	}
}

func (ctx *Context) GetNodeId() uint32 {
	return ctx.nodeId
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

func (ctx *Context) GetClients() []transport.Client {
	return ctx.clients
}

func (ctx *Context) BecomeFollower() {
	ctx.setType(domain.FOLLOWER)
	// TODO: do something with tickers
}

func (ctx *Context) BecomeCandidate() {
	ctx.setType(domain.CANDIDATE)
	// TODO: do something with tickers
}

func (ctx *Context) BecomeLeader() {
	ctx.setType(domain.LEADER)
	// TODO: do something with tickers}
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

func (ctx *Context) setType(target int) {
	ctx.nodeTypeMutex.Lock()
	ctx.nodeType = target
	ctx.nodeTypeMutex.Unlock()
}

func (ctx *Context) hasType(target int) bool {
	ctx.nodeTypeMutex.Lock()
	result := ctx.nodeType == target
	ctx.nodeTypeMutex.Unlock()
	return result
}
