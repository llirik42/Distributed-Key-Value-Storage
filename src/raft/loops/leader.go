package loops

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/raft/utils"
)

func LeaderLoop(ctx *context.Context) {
	ticker := ctx.GetLeaderLoopTicker()

	for range ticker.C {
		utils.SendHeartbeat(ctx)
	}
}
