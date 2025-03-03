package loops

import (
	"distributed-algorithms/raft/dto"
	"distributed-algorithms/raft/node"
	"time"
)

func LeaderLoop(node *node.Context, ticker *time.Ticker) {
	nodeId := node.GetId()

	for range ticker.C {
		request := dto.AppendEntriesRequest{
			Term:         node.GetCurrentTerm(),
			LeaderId:     int32(nodeId),
			PrevLogIndex: 0, // TODO
			PrevLogTerm:  0, // TODO
			LeaderCommit: 0, // TODO
		}

		// Sending heartbeat
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

func handleAppendEntriesResponse(node *node.Context, response *dto.AppendEntriesResponse) {
	// TODO: add checks related to logs

	node.CheckTerm(response.Term)
}
