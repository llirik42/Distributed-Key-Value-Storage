package context

import (
	"distributed-algorithms/src/config"
	"distributed-algorithms/src/log"
	"distributed-algorithms/src/raft/transport"
	"math/rand"
	"sync"
	"time"
)

type Context struct {
	ctxMutex sync.Mutex

	nodeId   string
	nodeRole int

	currentTerm uint32

	commitIndex uint64
	lastApplied uint64
	nextIndex   []uint64
	matchIndex  []uint64

	voted      bool
	votedFor   string
	voteNumber uint32

	log                         log.Log
	followerCandidateLoopTicker *time.Ticker
	leaderLoopTicker            *time.Ticker
	cfg                         config.RaftConfig
	server                      *transport.Server
	clients                     []transport.Client
}

func NewContext(cfg config.RaftConfig) *Context {
	ctx := &Context{
		nodeId:      cfg.SelfNode.Id,
		nodeRole:    Follower,
		currentTerm: 0,
		voted:       false,
		votedFor:    "", // Default value doesn't matter because voted = false by default
		voteNumber:  0,

		followerCandidateLoopTicker: nil,
		leaderLoopTicker:            nil,
		cfg:                         cfg,
		server:                      nil,
		clients:                     nil,
	}

	return ctx
}

func (ctx *Context) Lock() {
	ctx.ctxMutex.Lock()
}

func (ctx *Context) Unlock() {
	ctx.ctxMutex.Unlock()
}

func (ctx *Context) SetLog(log log.Log) {
	ctx.log = log
}

func (ctx *Context) SetServer(server *transport.Server) {
	ctx.server = server
}

func (ctx *Context) SetClients(clients []transport.Client) {
	ctx.clients = clients
}

func (ctx *Context) StartTickers() {
	ctx.followerCandidateLoopTicker = time.NewTicker(getRandomElectionTimeout(&ctx.cfg))
	ctx.leaderLoopTicker = time.NewTicker(getBroadcastTimeout(&ctx.cfg))
}

func (ctx *Context) GetFollowerCandidateLoopTicker() *time.Ticker {
	return ctx.followerCandidateLoopTicker
}

func (ctx *Context) GetLeaderLoopTicker() *time.Ticker {
	return ctx.leaderLoopTicker
}

func (ctx *Context) ResetNewElectionTimeout() {
	ctx.followerCandidateLoopTicker.Reset(getRandomElectionTimeout(&ctx.cfg))
}

func (ctx *Context) GetLeaderId() string {
	return "123" // TODO
}

func (ctx *Context) GetLogEntryTerm(index uint64) (uint32, bool) {
	term, exists, err := ctx.log.GetLogEntryTerm(index)

	if err != nil {
		// TODO: handle error (panic?)
		return 0, false
	}

	return term, exists
}

func (ctx *Context) PushCommand(cmd *log.Command) error {

	return nil
}

func (ctx *Context) GetLastLogEntryMetadata() log.EntryMetadata {
	metadata, err := ctx.log.GetLastLogEntryMetadata()

	if err != nil {
		// TODO: handle error (panic?)
		return log.EntryMetadata{}
	}

	return metadata
}

func (ctx *Context) ApplyByCommitIndex() {

}

func (ctx *Context) GetCommitIndex() uint64 {
	return ctx.commitIndex
}

func (ctx *Context) GetClusterSize() uint32 {
	return uint32(1 + len(ctx.cfg.OtherNodes))
}

func (ctx *Context) GetClients() []transport.Client {
	return ctx.clients
}

func (ctx *Context) GetNodeId() string {
	return ctx.nodeId
}

func (ctx *Context) IsFollower() bool {
	return ctx.hasRole(Follower)
}

func (ctx *Context) IsCandidate() bool {
	return ctx.hasRole(Candidate)
}

func (ctx *Context) IsLeader() bool {
	return ctx.hasRole(Leader)
}

func (ctx *Context) setRole(target int) {
	ctx.nodeRole = target
}

func (ctx *Context) hasRole(target int) bool {
	return ctx.nodeRole == target
}

func (ctx *Context) Vote(candidateId string) bool {
	var result bool

	if !ctx.voted {
		ctx.votedFor = candidateId
		ctx.voted = true
		result = true

		if candidateId == ctx.nodeId {
			// Node is a candidate and successfully votes for itself
			// Vote number must be persistent. If it's not, this condition must be added to case "ctx.voted"
			ctx.IncrementVoteNumber()
		}
	} else {
		// If the candidate we voted for falls, we'll again vote for it
		// If result = false, in this case we won't vote for the candidate
		result = ctx.votedFor == candidateId
	}

	return result
}

func (ctx *Context) ResetVoteNumber() {
	ctx.voteNumber = 0
}

func (ctx *Context) IncrementVoteNumber() uint32 {
	ctx.voteNumber++
	return ctx.voteNumber
}

func (ctx *Context) SetCurrentTerm(value uint32) {
	ctx.resetVoted()
	ctx.currentTerm = value
}

func (ctx *Context) IncrementCurrentTerm() uint32 {
	ctx.resetVoted()
	ctx.currentTerm++
	return ctx.currentTerm
}

func (ctx *Context) GetCurrentTerm() uint32 {
	return ctx.currentTerm
}

func (ctx *Context) BecomeFollower() {
	ctx.setRole(Follower)
	ctx.leaderLoopTicker.Stop()
	ctx.ResetNewElectionTimeout()
}

func (ctx *Context) BecomeCandidate() {
	ctx.setRole(Candidate)
	ctx.leaderLoopTicker.Stop()
	ctx.ResetNewElectionTimeout()
}

func (ctx *Context) BecomeLeader() {
	ctx.setRole(Leader)
	ctx.followerCandidateLoopTicker.Stop()
	ctx.leaderLoopTicker.Reset(getBroadcastTimeout(&ctx.cfg))
}

func (ctx *Context) resetVoted() {
	ctx.voted = false
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
