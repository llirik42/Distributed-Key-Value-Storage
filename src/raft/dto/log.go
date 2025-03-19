package dto

const (
	Set = iota
	Delete
)

type Command struct {
	Key   string
	Value any
	Type  int
}

type LogEntry struct {
	Index   uint64
	Term    uint32
	Command Command
}
