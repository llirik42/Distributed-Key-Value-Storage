package loops

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/raft/domain"
)

func FollowerCandidateLoop(ctx *context.Context) {
	ticker := ctx.GetFollowerCandidateLoopTicker()

	for range ticker.C {
		iterateFollowerCandidateLoop(ctx)
	}
}

func iterateFollowerCandidateLoop(ctx *context.Context) {
	ctx.Lock()
	defer ctx.Unlock()

	if ctx.IsFollower() {
		ctx.BecomeCandidate()
	}

	startNewTerm(ctx)
}

func startNewTerm(ctx *context.Context) {
	currentTerm := ctx.IncrementCurrentTerm()
	ctx.ResetVoteNumber()
	ctx.Vote(ctx.GetNodeId()) // Node votes for itself
	offerCandidacy(ctx, currentTerm)
}

func offerCandidacy(ctx *context.Context, currentTerm uint32) {
	request := domain.RequestVoteRequest{
		Term:         currentTerm,
		CandidateId:  ctx.GetNodeId(),
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	for _, client := range ctx.GetClients() {
		go func() {
			if err := client.SendRequestForVote(request); err != nil {
				// TODO: handle error
			}
		}()
	}
}
