package loops

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/utils"
	"time"
)

func FollowerCandidateLoop(ctx *context.Context, ticker *time.Ticker) {
	for range ticker.C {
		if ctx.IsFollower() {
			ctx.BecomeCandidate()
			startNewTerm(ctx)
		} else if ctx.IsCandidate() {
			clusterSize := ctx.GetClusterSize()

			if ctx.GetVoteNumber() > clusterSize/2 {
				ctx.BecomeLeader()
				sendHeartbeat(ctx)
			} else {
				ctx.SetNewRandomElectionTimeout()
				startNewTerm(ctx)
			}
		}
	}
}

func startNewTerm(ctx *context.Context) {
	currentTerm := ctx.IncrementCurrentTerm()
	offerCandidacy(ctx, currentTerm)
	ctx.ResetVoteNumber()
	ctx.Vote(ctx.GetNodeId()) // Node votes for itself
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
			response, err := client.SendRequestForVote(request)

			if err != nil {
				// TODO: handle error
			} else {
				handleRequestForVoteResponse(ctx, response)
			}
		}()
	}
}

func handleRequestForVoteResponse(ctx *context.Context, response *domain.RequestVoteResponse) {
	utils.CheckTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor

	if response.VoteGranted {
		ctx.IncrementVoteNumber()
	}
}
