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
			// TODO: становиться лидером и слать heartbeat нужно не здесь, а сразу при получении очередного положительного ответа на ReqeustForVote
			clusterSize := ctx.GetClusterSize()
			voteNumber := int(ctx.GetVoteNumber())

			if voteNumber > clusterSize/2 {
				ctx.BecomeLeader()
				sendHeartbeat(ctx)
			} else {
				startNewTerm(ctx)
			}
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
