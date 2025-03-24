package log

type Log interface {
	GetLogEntryTerm(index uint64) uint32

	TryGetLogEntryTerm(index uint64) (uint32, bool)

	GetLastLogEntryMetadata() EntryMetadata

	AddLogEntry(entry *Entry, index uint64)

	GetLastIndex() uint64

	GetLogEntries(startingIndex uint64) []Entry
}
