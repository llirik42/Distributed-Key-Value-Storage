package log

const (
	Set = iota
	Delete
)

type Command struct {
	Key   string
	Value any
	Type  int
}

type Entry struct {
	Term    uint32
	Command Command
}
