package utils

import (
	"distributed-key-value-storage/src/context"
	"distributed-key-value-storage/src/raft/dto"
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
		entries := logStorage.GetEntries(nextIndex)

		lastSentIndex := prevLogIndex + uint64(len(entries))
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
