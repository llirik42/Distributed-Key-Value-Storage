package log

import (
	"distributed-algorithms/src/key-value"
	"fmt"
)

func applyCommand(cmd *Command, storage key_value.Storage) error {
	switch cmd.Type {
	case Set:
		if err := storage.Set(cmd.Key, cmd.Value); err != nil {
			return fmt.Errorf("setting value %v of key %v failed: %v", cmd.Value, cmd.Key, err)
		}
	case Delete:
		if err := storage.Delete(cmd.Key); err != nil {
			return fmt.Errorf("deleting key %v failed: %v", cmd.Key, err)
		}
	default:
		return fmt.Errorf("unknown type of comman %v", cmd)
	}

	return nil
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
