package context

import (
	"distributed-key-value-storage/src/config"
	"distributed-key-value-storage/src/key-value"
	"distributed-key-value-storage/src/log"
	"distributed-key-value-storage/src/raft/transport"
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

type Context struct {
	ctxMutex sync.Mutex

	nodeId   string
	nodeRole int
	leaderId string

	currentTerm uint32

	commitIndex   uint64
	lastApplied   uint64
	nextIndex     []uint64
	matchIndex    []uint64
	lastSentIndex []uint64

	voted      bool
	votedFor   string
	voteNumber uint32

	followerCandidateLoopTicker *time.Ticker
	leaderLoopTicker            *time.Ticker
	cfg                         config.RaftConfig
	logStorage                  log.Storage
	commandExecutor             log.CommandExecutor
	keyValueStorage             key_value.Storage
	clients                     []transport.Client
}

func NewContext(cfg config.RaftConfig) *Context {
	otherNodesCount := len(cfg.OtherNodes)

	ctx := &Context{
		ctxMutex: sync.Mutex{},
		nodeId:   cfg.SelfNode.Id,
		nodeRole: 0, // Default value doesn't matter,
		// because BecomeFollower should be called before using context
		currentTerm:                 0,
		commitIndex:                 0,
		lastApplied:                 0,
		nextIndex:                   make([]uint64, otherNodesCount),
		matchIndex:                  make([]uint64, otherNodesCount),
		lastSentIndex:               make([]uint64, otherNodesCount),
		voted:                       false,
		votedFor:                    "", // Default value doesn't matter, because voted = false
		voteNumber:                  0,
		followerCandidateLoopTicker: nil,
		leaderLoopTicker:            nil,
		cfg:                         cfg,
		logStorage:                  nil,
		keyValueStorage:             nil,
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

func (ctx *Context) GetLogStorage() log.Storage {
	return ctx.logStorage
}

func (ctx *Context) SetLogStorage(storage log.Storage) {
	ctx.logStorage = storage
}

func (ctx *Context) GetCommandExecutor() log.CommandExecutor {
	return ctx.commandExecutor
}

func (ctx *Context) SetCommandExecutor(executor log.CommandExecutor) {
	ctx.commandExecutor = executor
}

func (ctx *Context) GetKeyValueStorage() key_value.Storage {
	return ctx.keyValueStorage
}

func (ctx *Context) SetKeyValueStorage(storage key_value.Storage) {
	ctx.keyValueStorage = storage
}

func (ctx *Context) GetClients() []transport.Client {
	return ctx.clients
}

func (ctx *Context) SetClients(clients []transport.Client) {
	ctx.clients = clients
}

func (ctx *Context) GetNodeId() string {
	return ctx.nodeId
}

func (ctx *Context) GetFollowerCandidateLoopTicker() *time.Ticker {
	return ctx.followerCandidateLoopTicker
}

func (ctx *Context) GetLeaderLoopTicker() *time.Ticker {
	return ctx.leaderLoopTicker
}

func (ctx *Context) GetLeaderId() string {
	return ctx.leaderId
}

func (ctx *Context) SetLeaderId(value string) {
	ctx.leaderId = value
}

func (ctx *Context) GetCommitIndex() uint64 {
	return ctx.commitIndex
}

func (ctx *Context) SetCommitIndex(value uint64) {
	ctx.commitIndex = value
	ctx.applyCommitedEntries()
}

func (ctx *Context) GetLastApplied() uint64 {
	return ctx.lastApplied
}

func (ctx *Context) GetNextIndex(clientIndex int) uint64 {
	return ctx.nextIndex[clientIndex]
}

func (ctx *Context) GetNextIndexes() []uint64 {
	return ctx.nextIndex
}

func (ctx *Context) SetNextIndex(clientIndex int, value uint64) {
	ctx.nextIndex[clientIndex] = value
}

func (ctx *Context) DecrementNextIndex(clientIndex int) {
	ctx.nextIndex[clientIndex]--
}

func (ctx *Context) GetLastSentIndex(clientIndex int) uint64 {
	return ctx.lastSentIndex[clientIndex]
}

func (ctx *Context) SetLastSentIndex(clientIndex int, value uint64) {
	ctx.lastSentIndex[clientIndex] = value
}

func (ctx *Context) GetMatchIndexes() []uint64 {
	return ctx.matchIndex
}

func (ctx *Context) SetMatchIndex(clientIndex int, value uint64) {
	ctx.matchIndex[clientIndex] = value

	newCommitIndex := ctx.findNewCommitIndex()
	if newCommitIndex != ctx.commitIndex {
		ctx.SetCommitIndex(newCommitIndex)
	}
}

func (ctx *Context) GetClusterSize() uint32 {
	return 1 + uint32(len(ctx.cfg.OtherNodes))
}

func (ctx *Context) PushCommand(cmd log.Command) string {
	commandId := uuid.NewString()

	cmd.Id = commandId

	entry := log.Entry{
		Term:    ctx.currentTerm,
		Command: cmd,
	}

	ctx.logStorage.PushLogEntry(entry)

	return commandId
}

func (ctx *Context) StartTickers() {
	ctx.followerCandidateLoopTicker = time.NewTicker(getRandomElectionTimeout(&ctx.cfg))
	ctx.leaderLoopTicker = time.NewTicker(getBroadcastTimeout(&ctx.cfg))
}

func (ctx *Context) ResetElectionTimeout() {
	ctx.followerCandidateLoopTicker.Reset(getRandomElectionTimeout(&ctx.cfg))
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
	ctx.resetLeaderId()
	ctx.currentTerm = value
}

func (ctx *Context) IncrementCurrentTerm() uint32 {
	ctx.resetVoted()
	ctx.resetLeaderId()
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
	ctx.leaderId = ctx.nodeId
	ctx.followerCandidateLoopTicker.Stop()
	ctx.leaderLoopTicker.Reset(getBroadcastTimeout(&ctx.cfg))
	ctx.initNextAndMatchIndexes()
}

func (ctx *Context) applyCommitedEntries() {
	for i := ctx.lastApplied + 1; i <= ctx.commitIndex; i++ {
		cmd := ctx.logStorage.GetEntryCommand(i)
		ctx.commandExecutor.Execute(cmd)
		ctx.lastApplied++
	}
}

func (ctx *Context) initNextAndMatchIndexes() {
	for i := 0; i < len(ctx.clients); i++ {
		ctx.nextIndex[i] = 1 // TODO: modify that leader won't send ALL entries at start
		ctx.matchIndex[i] = 0
	}
}

func (ctx *Context) findNewCommitIndex() uint64 {
	currentTerm := ctx.currentTerm
	clusterSizeHalved := ctx.GetClusterSize() / 2

	logLength := ctx.logStorage.GetLength()

	newCommitIndex := ctx.commitIndex

	for i := ctx.commitIndex + 1; i <= logLength; i++ {
		var count uint32 = 1 // By default, include current node (leader)
		for _, v := range ctx.matchIndex {
			if v >= i {
				count++
			}
		}

		if count <= clusterSizeHalved {
			break
		}

		term := ctx.logStorage.GetEntryMetadata(i).Term
		if term != currentTerm {
			continue
		}

		newCommitIndex = i
	}

	return newCommitIndex
}

func (ctx *Context) resetVoted() {
	ctx.voted = false
}

func (ctx *Context) resetLeaderId() {
	ctx.leaderId = ""
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
