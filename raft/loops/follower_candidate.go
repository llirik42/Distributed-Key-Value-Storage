package loops

import (
	"distributed-algorithms/raft/dto"
	"distributed-algorithms/raft/node"
	"math/rand"
	"time"
)

func FollowerCandidateLoop(node *node.Node, ticker *time.Ticker) {
	minElectionTimeoutMs := 500
	maxElectionTimeoutMs := 700
	duration := getElectionTimeout(minElectionTimeoutMs, maxElectionTimeoutMs)

	for range ticker.C {
		if node.IsFollower() {
			node.BecomeCandidate()

			// Начать голосование
			node.ResetVoteNumber()
			currentTerm := node.IncrementCurrentTerm()
			request := dto.RequestVoteRequest{
				Term:         currentTerm,
				CandidateId:  int32(node.GetId()),
				LastLogIndex: 0,
				LastLogTerm:  0,
			}

			// Sending requests for vote
			for _, client := range node.GetClients() {
				go func() {
					response, err := client.SendRequestForVote(request)

					if err != nil {
						// TODO: handle error
					} else {
						handleRequestForVoteResponse(node, response)
					}
				}()
			}
		} else if node.IsCandidate() {
			// TODO: read from config
			var nodeCount int32 = 5

			if node.GetVoteNumber() > nodeCount/2 {
				// become leader and immediately send heartbeat
			} else {
				currentTerm := node.IncrementCurrentTerm()

				request := dto.RequestVoteRequest{
					Term:         currentTerm,
					CandidateId:  int32(node.GetId()),
					LastLogIndex: 0,
					LastLogTerm:  0,
				}

				// Sending requests for vote
				for _, client := range node.GetClients() {
					go func() {
						response, err := client.SendRequestForVote(request)

						if err != nil {
							// TODO: handle error
						} else {
							handleRequestForVoteResponse(node, response)
						}
					}()
				}
			}
		}
	}
}

func handleRequestForVoteResponse(node *node.Node, response *dto.RequestVoteResponse) {
	node.CheckTerm(response.Term) // TODO: Check this in gRPC-interceptor

	if response.VoteGranted {
		node.IncrementVoteNumber()
	}
}

func getElectionTimeout(minElectionTimeoutMs int, maxElectionTimeoutMs int) time.Duration {
	electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs
	return getDurationMs(electionTimeoutMs)
}
