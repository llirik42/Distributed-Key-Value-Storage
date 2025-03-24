package log

import (
	"distributed-algorithms/src/key-value"
	"fmt"
)

func ApplyCommand(cmd *Command, storage key_value.Storage) {
	switch cmd.Type {
	case Set:
		storage.Set(cmd.Key, cmd.Value)
	case Delete:
		storage.Delete(cmd.Key)
	default:
		panic(fmt.Errorf("unknown type of command %+v", *cmd))
	}
}

// CompareEntries Compares two log entries by term and index
//
// Return value:
//
// -1, if first log is more up-to-date
//
// 0, if no log is more up-to-date
//
// 1, if second log is more up-to-date
func CompareEntries(term1 uint32, index1 uint64, term2 uint32, index2 uint64) int {
	if term1 < term2 {
		return 1
	}
	if term1 > term2 {
		return -1
	}

	if index1 < index2 {
		return 1
	}
	if index1 > index2 {
		return -1
	}

	return 0
}
