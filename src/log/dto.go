package log

const (
	Set = iota
	Delete
)

type EntryMetadata struct {
	Term  uint32
	Index uint64
}

type Command struct {
	Key   string
	Value any
	Type  int
}

type Entry struct {
	Term    uint32
	Command Command
}
