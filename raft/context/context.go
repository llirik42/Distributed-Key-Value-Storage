package context

import (
	"distributed-algorithms/raft"
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

type Context struct {
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

func (ctx *Context) Start() error {
	ctx.nodeType = follower // is it necessary to use ctx.setType?

	return ctx.server.Listen()
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

func (ctx *Context) GetId() Id {
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
	ctx.setType(follower)
	// TODO: do something with tickers
}

func (ctx *Context) BecomeCandidate() {
	ctx.setType(candidate)
	// TODO: do something with tickers
}

func (ctx *Context) BecomeLeader() {
	ctx.setType(leader)
	// TODO: do something with tickers}
}

func (ctx *Context) IsFollower() bool {
	return ctx.hasType(follower)
}

func (ctx *Context) IsCandidate() bool {
	return ctx.hasType(candidate)
}

func (ctx *Context) IsLeader() bool {
	return ctx.hasType(leader)
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
