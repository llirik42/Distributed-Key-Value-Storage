package context

import (
	"distributed-algorithms/config"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/transport"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Context struct {
	cfg config.RaftConfig

	currentTerm atomic.Int32

	voted         bool
	votedFor      string
	votedForMutex sync.Mutex
	voteNumber    atomic.Uint32

	nodeId        string
	nodeRole      int
	nodeRoleMutex sync.Mutex

	followerCandidateLoopTicker *time.Ticker
	leaderLoopTicker            *time.Ticker

	server  *transport.Server
	clients []transport.Client
}

func NewContext(cfg config.RaftConfig) *Context {
	ctx := &Context{
		cfg:                         cfg,
		currentTerm:                 atomic.Int32{},
		voted:                       false,
		votedFor:                    "",
		votedForMutex:               sync.Mutex{},
		voteNumber:                  atomic.Uint32{},
		nodeId:                      cfg.SelfNode.Id,
		nodeRole:                    domain.FOLLOWER,
		nodeRoleMutex:               sync.Mutex{},
		followerCandidateLoopTicker: nil,
		leaderLoopTicker:            nil,
		server:                      nil,
		clients:                     nil,
	}

	return ctx
}

func (ctx *Context) SetServer(server *transport.Server) {
	ctx.server = server
}

func (ctx *Context) SetClients(clients []transport.Client) {
	ctx.clients = clients
}

func (ctx *Context) StartTickers() {
	followerCandidateLoopTicker := time.NewTicker(getRandomElectionTimeout(&ctx.cfg))
	leaderLoopTicker := time.NewTicker(getBroadcastTimeout(&ctx.cfg))

	ctx.followerCandidateLoopTicker = followerCandidateLoopTicker
	ctx.leaderLoopTicker = leaderLoopTicker
}

func (ctx *Context) GetFollowerCandidateLoopTicker() *time.Ticker {
	return ctx.followerCandidateLoopTicker
}

func (ctx *Context) GetLeaderLoopTicker() *time.Ticker {
	return ctx.leaderLoopTicker
}

func (ctx *Context) SetNewRandomElectionTimeout() {
	ctx.followerCandidateLoopTicker.Reset(getRandomElectionTimeout(&ctx.cfg))
}

func (ctx *Context) GetClusterSize() int {
	return 1 + len(ctx.cfg.OtherNodes)
}

func (ctx *Context) GetClients() []transport.Client {
	return ctx.clients
}

func (ctx *Context) GetNodeId() string {
	return ctx.nodeId
}

func (ctx *Context) IsFollower() bool {
	return ctx.hasRole(domain.FOLLOWER)
}

func (ctx *Context) IsCandidate() bool {
	return ctx.hasRole(domain.CANDIDATE)
}

func (ctx *Context) IsLeader() bool {
	return ctx.hasRole(domain.LEADER)
}

func (ctx *Context) setRole(target int) {
	ctx.nodeRoleMutex.Lock()
	defer ctx.nodeRoleMutex.Unlock()
	ctx.nodeRole = target
}

func (ctx *Context) hasRole(target int) bool {
	ctx.nodeRoleMutex.Lock()
	defer ctx.nodeRoleMutex.Unlock()
	result := ctx.nodeRole == target
	return result
}

func (ctx *Context) ResetVotedFor() {
	ctx.votedForMutex.Lock()
	defer ctx.votedForMutex.Unlock()
	ctx.voted = false
}

func (ctx *Context) Vote(candidateId string) bool {
	ctx.votedForMutex.Lock()
	defer ctx.votedForMutex.Unlock()

	// Whether node votes for itself
	if ctx.nodeId == candidateId {
		ctx.votedFor = candidateId
		ctx.voted = true
		ctx.IncrementVoteNumber()
		return true
	}

	var result bool

	if !ctx.voted {
		ctx.votedFor = candidateId
		ctx.voted = true
		result = true
	} else {
		// If the candidate we voted for falls, we'll again vote for it
		// If result = false, in this case we won't vote for the candidate
		result = ctx.votedFor == candidateId
	}

	return result
}

func (ctx *Context) ResetVoteNumber() {
	ctx.voteNumber.Store(0)
}

func (ctx *Context) IncrementVoteNumber() {
	ctx.voteNumber.Add(1)
}

func (ctx *Context) GetVoteNumber() uint32 {
	return ctx.voteNumber.Load()
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
	ctx.leaderLoopTicker.Stop()
	ctx.SetNewRandomElectionTimeout()
}

func (ctx *Context) BecomeCandidate() {
	ctx.setRole(domain.CANDIDATE)
	ctx.leaderLoopTicker.Stop()
}

func (ctx *Context) BecomeLeader() {
	ctx.setRole(domain.LEADER)
	ctx.followerCandidateLoopTicker.Stop()
	ctx.leaderLoopTicker.Reset(getBroadcastTimeout(&ctx.cfg))
}

func getBroadcastTimeout(cfg *config.RaftConfig) time.Duration {
	return getDurationMs(cfg.BroadcastTimeMs)
}

func getRandomElectionTimeout(cfg *config.RaftConfig) time.Duration {
	electionTimeoutMs := rand.Intn(cfg.MaxElectionTimeoutMs-cfg.MinElectionTimeoutMs+1) + cfg.MinElectionTimeoutMs
	return getDurationMs(electionTimeoutMs)
}

func getDurationMs(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}
