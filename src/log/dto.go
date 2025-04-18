package log

// Type of command
const (
	Get = iota
	Set
	CompareAndSet
	Delete
	AddElement
)

type CommandExecutionInfo struct {
	Value   any
	Message string
	Success bool
}

type EntryMetadata struct {
	Term  uint32
	Index uint64
}

type Command struct {
	Id       string
	Key      string
	SubKey   string
	OldValue any
	NewValue any
	Type     int
}

type Entry struct {
	Term    uint32
	Command Command
}
