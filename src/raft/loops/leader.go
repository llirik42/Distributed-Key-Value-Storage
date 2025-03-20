package loops

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/raft/utils"
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
	utils.SendHeartbeat(ctx)
}
