package log

type Storage interface {
	GetEntryMetadata(index uint64) EntryMetadata

	GetEntryCommand(index uint64) Command

	GetLastEntryMetadata() EntryMetadata

	GetLogEntries(startIndex uint64) []Entry

	TryGetEntryMetadata(index uint64) (EntryMetadata, bool)

	PushLogEntry(entry Entry)

	AddLogEntry(entry Entry, index uint64)
}
