package loops

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/domain"
	"math/rand"
	"time"
)

func FollowerCandidateLoop(ctx *context.Context, ticker *time.Ticker) {
	minElectionTimeoutMs := 500
	maxElectionTimeoutMs := 700

	// TODO: избавиться от копипасты
	for range ticker.C {
		if ctx.IsFollower() {
			ctx.BecomeCandidate()

			// Начать голосование

			// TODO: изменить votedFor на свой Id
			ctx.ResetVoteNumber()
			currentTerm := ctx.IncrementCurrentTerm()
			request := domain.RequestVoteRequest{
				Term:         currentTerm,
				CandidateId:  ctx.GetNodeId(),
				LastLogIndex: 0,
				LastLogTerm:  0,
			}

			// Sending requests for vote
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
		} else if ctx.IsCandidate() {
			// TODO: read from config
			var nodeCount int32 = 5

			if ctx.GetVoteNumber() > nodeCount/2 {
				// become leader and immediately send heartbeat
			} else {
				// Choose new election timeout
				ticker.Reset(GetElectionTimeout(minElectionTimeoutMs, maxElectionTimeoutMs))

				// Reset ticker to new random timeout
				// TODO: изменить votedFor на свой Id
				ctx.ResetVoteNumber()
				currentTerm := ctx.IncrementCurrentTerm()
				request := domain.RequestVoteRequest{
					Term:         currentTerm,
					CandidateId:  ctx.GetNodeId(),
					LastLogIndex: 0,
					LastLogTerm:  0,
				}

				// Sending requests for vote
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
		}
	}
}

func handleRequestForVoteResponse(ctx *context.Context, response *domain.RequestVoteResponse) {
	ctx.CheckTerm(response.Term) // TODO: Check this in gRPC-interceptor

	if response.VoteGranted {
		ctx.IncrementVoteNumber()
	}
}

func GetElectionTimeout(minElectionTimeoutMs int, maxElectionTimeoutMs int) time.Duration {
	electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs
	return getDurationMs(electionTimeoutMs)
}
