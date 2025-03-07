package loops

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/utils"
)

func LeaderLoop(ctx *context.Context) {
	ticker := ctx.GetLeaderLoopTicker()

	for range ticker.C {
		utils.SendHeartbeat(ctx)
	}
}
