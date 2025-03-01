package loops

import (
	"distributed-algorithms/raft/node"
	"time"
)

func LeaderLoop(node *node.Node) {
	broadcastTimeoutMs := 100
	duration := getDurationMs(broadcastTimeoutMs)

	for {
		time.Sleep(duration)

		// Node stopped being leader
		if !node.IsLeader() {
			break
		}

		// Send heartbeat

	}
}
