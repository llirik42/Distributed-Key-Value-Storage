package log

type Storage interface {
	GetLength() uint64

	FindFirstEntryWithTerm(term uint32) (uint64, bool)

	FindLastEntryWithTerm(term uint32) (uint64, bool)

	GetEntryMetadata(index uint64) EntryMetadata

	GetEntryCommand(index uint64) Command

	GetLastEntryMetadata() EntryMetadata

	GetEntries(startIndex uint64) []Entry

	TryGetEntryMetadata(index uint64) (EntryMetadata, bool)

	PushLogEntry(entry Entry)

	AddLogEntry(entry Entry, index uint64)
}

type CommandExecutor interface {
	Execute(cmd Command)

	GetCommandExecutionInfo(commandId string) (CommandExecutionInfo, bool)
}
