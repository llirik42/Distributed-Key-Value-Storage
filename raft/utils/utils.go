package utils

import (
	"distributed-algorithms/context"
	"distributed-algorithms/raft/domain"
)

func SendHeartbeat(ctx *context.Context) {
	ctx.Lock()
	defer ctx.Unlock()

	request := domain.AppendEntriesRequest{
		Term:         ctx.GetCurrentTerm(),
		LeaderId:     ctx.GetNodeId(),
		PrevLogIndex: 0, // TODO
		PrevLogTerm:  0, // TODO
		LeaderCommit: 0, // TODO
	}

	for _, client := range ctx.GetClients() {
		go func() {
			if err := client.SendAppendEntries(request); err != nil {
				// TODO: handle error
			}
		}()
	}
}
