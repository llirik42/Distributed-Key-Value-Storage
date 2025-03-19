package domain

const (
	Set = iota
	Delete
)

type Command struct {
	Key   string
	Value any
	Type  int
}
