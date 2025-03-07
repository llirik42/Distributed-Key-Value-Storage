package loops

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/domain"
)

func FollowerCandidateLoop(ctx *context.Context) {
	ticker := ctx.GetFollowerCandidateLoopTicker()

	for range ticker.C {
		if ctx.IsFollower() {
			ctx.BecomeCandidate()
			startNewTerm(ctx)
		} else if ctx.IsCandidate() {
			// The candidate didn't get quorum
			startNewTerm(ctx)
		}
	}
}

func startNewTerm(ctx *context.Context) {
	currentTerm := ctx.IncrementCurrentTerm()

	ctx.ResetVoteNumber()
	ctx.Vote(ctx.GetNodeId()) // Node votes for itself
	offerCandidacy(ctx, currentTerm)
}

func offerCandidacy(ctx *context.Context, currentTerm int32) {
	request := domain.RequestVoteRequest{
		Term:         currentTerm,
		CandidateId:  ctx.GetNodeId(),
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	for _, client := range ctx.GetClients() {
		go func() {
			err := client.SendRequestForVote(request)

			if err != nil {
				// TODO: handle error
			}
		}()
	}
}
