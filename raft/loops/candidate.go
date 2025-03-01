package loops

import (
	"distributed-algorithms/raft/node"
	"math/rand"
	"time"
)

func CandidateLoop(node *node.Node) {
	minElectionTimeoutMs := 500
	maxElectionTimeoutMs := 700
	duration := getElectionTimeout(minElectionTimeoutMs, maxElectionTimeoutMs)

	for {
		time.Sleep(duration)

		// Node stopped being candidate
		if !node.IsCandidate() {
			break
		}

	}
}

func getElectionTimeout(minElectionTimeoutMs int, maxElectionTimeoutMs int) time.Duration {
	electionTimeoutMs := rand.Intn(maxElectionTimeoutMs-minElectionTimeoutMs+1) + minElectionTimeoutMs
	return getDurationMs(electionTimeoutMs)
}
