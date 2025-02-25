package raft

// CommonState TODO: separate persistent and volatile state
type CommonState struct {
	currentTerm int32 // Not nullable (init value = 0)
	votedFor    int32 // Nullable
	// TODO: add log[]
	commitIndex int // Not nullable
	lastApplied int // Not nullable
}

const noVotedFor = -1
