package loops

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/domain"
	"time"
)

func LeaderLoop(ctx *context.Context, ticker *time.Ticker) {
	nodeId := ctx.GetNodeId()

	for range ticker.C {
		request := domain.AppendEntriesRequest{
			Term:         ctx.GetCurrentTerm(),
			LeaderId:     nodeId,
			PrevLogIndex: 0, // TODO
			PrevLogTerm:  0, // TODO
			LeaderCommit: 0, // TODO
		}

		// Sending heartbeat
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
}

func handleAppendEntriesResponse(ctx *context.Context, response *domain.AppendEntriesResponse) {
	// TODO: add checks related to logs

	ctx.CheckTerm(response.Term)
}
