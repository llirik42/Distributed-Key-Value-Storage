package loops

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/utils"
	"encoding/json"
	"log"
)

func sendHeartbeat(ctx *context.Context) {
	request := domain.AppendEntriesRequest{
		Term:         ctx.GetCurrentTerm(),
		LeaderId:     ctx.GetNodeId(),
		PrevLogIndex: 0, // TODO
		PrevLogTerm:  0, // TODO
		LeaderCommit: 0, // TODO
	}

	for _, client := range ctx.GetClients() {
		go func() {
			response, err := client.SendAppendEntries(request)

			if err != nil {
				// TODO: handle error
			} else {
				handleAppendEntriesResponse(ctx, response)
			}
		}()
	}
}

func handleAppendEntriesResponse(ctx *context.Context, response *domain.AppendEntriesResponse) {
	// TODO: add checks related to logs

	a, _ := json.Marshal(response)
	log.Printf("Node \"%s\" received response of append-entries: %s", ctx.GetNodeId(), a)

	utils.CheckTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor
}
