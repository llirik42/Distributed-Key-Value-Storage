package loops

import (
	"distributed-algorithms/raft/context"
	"time"
)

func LeaderLoop(ctx *context.Context, ticker *time.Ticker) {
	for range ticker.C {
		sendHeartbeat(ctx)
	}
}
