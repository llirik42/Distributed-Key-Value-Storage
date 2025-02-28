package raft

import (
	"time"
)

func leaderLoop(node *Node) {
	broadcastTimeoutMs := 100
	duration := time.Duration(broadcastTimeoutMs) * time.Millisecond

	for {
		time.Sleep(duration)

		// Node stopped being leader
		if !node.IsLeader() {
			break
		}

		// Send heartbeat
		
	}
}
