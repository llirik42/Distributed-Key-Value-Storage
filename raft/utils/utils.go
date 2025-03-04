package utils

import "distributed-algorithms/raft/context"

func CheckTerm(ctx *context.Context, term int32) {
	if term > ctx.GetCurrentTerm() {

	}
}
