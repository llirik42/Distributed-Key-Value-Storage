package loops

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/domain"
	"distributed-algorithms/raft/utils"
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

	log.Printf("Node \"%s\" received response of append-entries: %v", ctx.GetNodeId(), response)

	utils.CheckTerm(ctx, response.Term) // TODO: Check this in gRPC-interceptor
}
