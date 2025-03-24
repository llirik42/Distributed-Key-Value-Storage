package in_memory

import (
	"distributed-algorithms/src/log"
	"fmt"
)

type Storage struct {
	entries  []log.Entry
	length   uint64
	capacity uint64
}

func NewStorage() log.Storage {
	var startCapacity uint64 = 1

	return &Storage{
		entries:  make([]log.Entry, startCapacity),
		length:   0,
		capacity: startCapacity,
	}
}

func (storage *Storage) GetEntryMetadata(index uint64) log.EntryMetadata {
	storage.validateIndex(index)

	if index == 0 {
		return log.EntryMetadata{}
	}

	entry := storage.getEntry(index)

	return log.EntryMetadata{
		Term:  entry.Term,
		Index: index,
	}
}

func (storage *Storage) GetEntryCommand(index uint64) log.Command {
	storage.validateIndex(index)

	if index == 0 {
		return log.Command{}
	}

	return storage.getEntry(index).Command
}

func (storage *Storage) GetLastEntryMetadata() log.EntryMetadata {
	if storage.isEmpty() {
		return log.EntryMetadata{}
	}

	entry := storage.getEntry(storage.length)

	return log.EntryMetadata{
		Term:  entry.Term,
		Index: storage.length,
	}
}

func (storage *Storage) GetLogEntries(startIndex uint64) []log.Entry {
	return storage.entries[getPhysicalIndex(startIndex):storage.length]
}

func (storage *Storage) TryGetEntryMetadata(index uint64) (log.EntryMetadata, bool) {
	if !storage.isIndexValid(index) {
		return log.EntryMetadata{}, false
	}

	if index == 0 {
		return log.EntryMetadata{}, true
	}

	return storage.GetEntryMetadata(index), true
}

func (storage *Storage) PushLogEntry(entry log.Entry) {
	// Check capacity
	if storage.length == storage.capacity {
		tmp := make([]log.Entry, storage.capacity)
		storage.entries = append(storage.entries, tmp...)
		storage.capacity *= 2
	}

	storage.entries[storage.length] = entry
	storage.length++
}

func (storage *Storage) AddLogEntry(entry log.Entry, index uint64) {
	if index > storage.length+1 || index == 0 {
		panic(fmt.Errorf("invalid index: %d", index))
	}

	if index == storage.length+1 {
		storage.PushLogEntry(entry)
	}

	storage.entries[index-1] = entry
	storage.length = min(storage.length, index)
}

func (storage *Storage) getEntry(index uint64) log.Entry {
	return storage.entries[getPhysicalIndex(index)]
}

func (storage *Storage) isEmpty() bool {
	return storage.length == 0
}

func (storage *Storage) validateIndex(index uint64) {
	if !storage.isIndexValid(index) {
		panic(fmt.Errorf("index %d out of range", index))
	}
}

func (storage *Storage) isIndexValid(index uint64) bool {
	return index <= storage.length
}

func getPhysicalIndex(index uint64) uint64 {
	return index - 1
}
