package log

type Log interface {
	AddLogEntry(entry *Entry, index uint64)

	PushLogEntry(entry *Entry)

	GetLogEntryTerm(index uint64) uint32

	GetLogEntry(index uint64) *Entry

	TryGetLogEntryTerm(index uint64) (uint32, bool)

	GetLastLogEntryMetadata() EntryMetadata

	GetLastIndex() uint64

	GetLogEntries(startingIndex uint64) []*Entry
}
