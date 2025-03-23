package log

type Log interface {
	GetLogEntryTerm(index uint64) (uint32, bool, error)

	GetLastLogEntryMetadata() (EntryMetadata, error)

	AddLogEntry(entry *Entry, index uint64) error

	GetLastIndex() (uint64, error)

	GetLogEntries(startingIndex uint64) ([]Entry, error)
}
