package loops

import (
	"distributed-algorithms/raft/context"
)

func LeaderLoop(ctx *context.Context) {
	ticker := ctx.GetLeaderLoopTicker()

	for range ticker.C {
		sendHeartbeat(ctx)
	}
}
