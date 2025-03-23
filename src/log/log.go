package log

type Log interface {
	GetLogEntryTerm(index uint64) (uint32, bool, error)

	GetLastLogEntryMetadata() (EntryMetadata, error)
}
