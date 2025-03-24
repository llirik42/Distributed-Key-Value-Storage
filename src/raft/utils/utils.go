package utils

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/raft/dto"
)

func SendAppendEntries(ctx *context.Context) {
	logStorage := ctx.GetLogStorage()
	term := ctx.GetCurrentTerm()
	leaderId := ctx.GetLeaderId()
	leaderCommit := ctx.GetCommitIndex()

	for _, client := range ctx.GetClients() {
		clientIndex := client.GetIndex()
		nextIndex := ctx.GetNextIndex(clientIndex)
		prevLogIndex := nextIndex - 1
		prevLogTerm := logStorage.GetEntryMetadata(prevLogIndex).Term
		entries := logStorage.GetLogEntries(nextIndex)

		lastSentIndex := prevLogIndex + uint64(len(entries)) - 1
		ctx.SetLastSentIndex(clientIndex, lastSentIndex)

		request := dto.AppendEntriesRequest{
			Term:         term,
			LeaderId:     leaderId,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			LeaderCommit: leaderCommit,
			Entries:      entries,
		}

		go func() {
			if err := client.SendAppendEntries(request); err != nil {
				// TODO: handle error
			}
		}()
	}
}
