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

	commitIndex   uint64
	lastApplied   uint64
	nextIndex     []uint64
	matchIndex    []uint64
	lastSentIndex []uint64

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
		ctxMutex:                    sync.Mutex{},
		nodeId:                      cfg.SelfNode.Id,
		nodeRole:                    Follower,
		currentTerm:                 0,
		commitIndex:                 0,
		lastApplied:                 0,
		nextIndex:                   []uint64{},
		matchIndex:                  []uint64{},
		lastSentIndex:               []uint64{},
		voted:                       false,
		votedFor:                    "",
		voteNumber:                  0,
		log:                         nil,
		followerCandidateLoopTicker: nil,
		leaderLoopTicker:            nil,
		cfg:                         config.RaftConfig{},
		server:                      nil,
		clients:                     nil,
	}

	// Init array-fields
	for i := 0; i < len(cfg.OtherNodes); i++ {
		ctx.nextIndex[i] = 0
		ctx.matchIndex[i] = 0
		ctx.lastSentIndex[i] = 0
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

func (ctx *Context) ResetElectionTimeout() {
	ctx.followerCandidateLoopTicker.Reset(getRandomElectionTimeout(&ctx.cfg))
}

func (ctx *Context) GetLeaderId() string {
	return "123" // TODO
}

func (ctx *Context) GetLastApplied() uint64 {
	return ctx.lastApplied
}

func (ctx *Context) GetNextIndexes() []uint64 {
	return ctx.nextIndex
}

func (ctx *Context) GetMatchIndexes() []uint64 {
	return ctx.matchIndex
}

func (ctx *Context) GetLog() log.Log {
	return ctx.log
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
	// TODO:

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

func (ctx *Context) SetCommitIndex(value uint64) {
	ctx.commitIndex = value

	// TODO: If commitIndex > lastApplied: increment lastApplied, apply
	//	log[lastApplied] to state machine (§5.3)

}

func (ctx *Context) GetCommitIndex() uint64 {
	return ctx.commitIndex
}

func (ctx *Context) AddLogEntry(entry *log.Entry, index uint64) {
	// TODO:
}

func (ctx *Context) GetLogEntries(startingIndex uint64) []log.Entry {
	entries, err := ctx.log.GetLogEntries(startingIndex)

	if err != nil {
		// TODO: handle error
		return []log.Entry{}
	}

	return entries
}

func (ctx *Context) SetNextIndex(clientIndex int, value uint64) {
	ctx.nextIndex[clientIndex] = value
}

func (ctx *Context) GetNextIndex(clientIndex int) uint64 {
	return ctx.nextIndex[clientIndex]
}

func (ctx *Context) DecrementNextIndex(clientIndex int) {
	ctx.nextIndex[clientIndex]--
}

func (ctx *Context) SetLastSentIndex(clientIndex int, value uint64) {
	ctx.lastSentIndex[clientIndex] = value
}

func (ctx *Context) GetLastSentIndex(clientIndex int) uint64 {
	return ctx.lastSentIndex[clientIndex]
}

func (ctx *Context) SetMatchIndex(clientIndex int, value uint64) {
	ctx.matchIndex[clientIndex] = value

	// Update commitIndex
	currentTerm := ctx.currentTerm
	clusterSizeHalved := ctx.GetClusterSize() / 2
	lastLogIndex := ctx.log.GetLastIndex()
	newCommitIndex := ctx.commitIndex
	for i := ctx.commitIndex + 1; i <= lastLogIndex; i++ {
		var count uint32 = 1 // By default, include current node (leader)
		for _, v := range ctx.matchIndex {
			if v >= i {
				count++
			}
		}

		if count <= clusterSizeHalved {
			break
		}

		term := ctx.log.GetLogEntryTerm(i)
		if term != currentTerm {
			continue
		}

		newCommitIndex = i
	}

	if newCommitIndex != ctx.commitIndex {
		ctx.SetCommitIndex(newCommitIndex)
	}
}

func (ctx *Context) GetClusterSize() uint32 {
	return 1 + uint32(len(ctx.cfg.OtherNodes))
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
	ctx.ResetElectionTimeout()
}

func (ctx *Context) BecomeCandidate() {
	ctx.setRole(Candidate)
	// Node can become Candidate only from Follower.
	// So at this point we don't have to deal with tickers, because Follower's tickers = Candidate's tickers
}

func (ctx *Context) BecomeLeader() {
	ctx.setRole(Leader)
	ctx.followerCandidateLoopTicker.Stop()
	ctx.leaderLoopTicker.Reset(getBroadcastTimeout(&ctx.cfg))
	ctx.initNextAndMatchIndexes()
}

func (ctx *Context) initNextAndMatchIndexes() {
	lastLogIndex, err := ctx.log.GetLastIndex()
	if err != nil {
		// TODO: handle error
	}

	for i := 0; i < len(ctx.clients); i++ {
		ctx.nextIndex[i] = lastLogIndex + 1
		ctx.matchIndex[i] = 0
	}
}

func (ctx *Context) resetVoted() {
	ctx.voted = false
}

func (ctx *Context) setRole(target int) {
	ctx.nodeRole = target
}

func (ctx *Context) hasRole(target int) bool {
	return ctx.nodeRole == target
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
