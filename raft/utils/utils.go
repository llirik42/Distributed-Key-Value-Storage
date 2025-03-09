package utils

import (
	"distributed-algorithms/context"
	"distributed-algorithms/raft/domain"
	"log"
)

func SendHeartbeat(ctx *context.Context) {
	request := domain.AppendEntriesRequest{
		Term:         ctx.GetCurrentTerm(),
		LeaderId:     ctx.GetNodeId(),
		PrevLogIndex: 0, // TODO
		PrevLogTerm:  0, // TODO
		LeaderCommit: 0, // TODO
	}

	log.Println("Sending append heartbeat")

	for _, client := range ctx.GetClients() {
		go func() {
			err := client.SendAppendEntries(request)

			if err != nil {
				// TODO: handle error
			}
		}()
	}
}
