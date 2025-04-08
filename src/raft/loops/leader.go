package loops

import (
	"distributed-key-value-storage/src/context"
	"distributed-key-value-storage/src/raft/utils"
)

func LeaderLoop(ctx *context.Context) {
	ticker := ctx.GetLeaderLoopTicker()

	for range ticker.C {
		iterateLeaderLoop(ctx)
	}
}

func iterateLeaderLoop(ctx *context.Context) {
	ctx.Lock()
	defer ctx.Unlock()
	utils.SendAppendEntries(ctx)
}
