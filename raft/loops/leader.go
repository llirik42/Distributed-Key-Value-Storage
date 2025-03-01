package loops

import (
	"distributed-algorithms/raft/dto"
	"distributed-algorithms/raft/node"
	"time"
)

func LeaderLoop(node *node.Node) {
	nodeId := node.GetId()

	broadcastTimeoutMs := 100 // TODO
	duration := getDurationMs(broadcastTimeoutMs)

	for {
		time.Sleep(duration)

		// Node stopped being leader
		if !node.IsLeader() {
			break
		}

		// Sending heartbeat
		request := dto.AppendEntriesRequest{
			Term:         node.GetCurrentTerm(),
			LeaderId:     nodeId,
			PrevLogIndex: 0, // TODO
			PrevLogTerm:  0, // TODO
			LeaderCommit: 0, // TODO
		}
		for _, client := range node.GetClients() {
			go func() {
				response, err := client.SendAppendEntries(request)

				if err != nil {
					// TODO: handle error
				} else {
					handleAppendEntriesResponse(node, response)
				}
			}()
		}
	}
}

func handleAppendEntriesResponse(node *node.Node, response *dto.AppendEntriesResponse) {
	// TODO: add checks related to logs

	responseTerm := response.Term

	if responseTerm > node.GetCurrentTerm() {
		node.SetCurrentTerm(responseTerm)
		node.BecomeFollower()
	}
}
