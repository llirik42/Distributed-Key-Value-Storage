package client_interaction

const (
	Get           = "get"
	Set           = "set"
	CompareAndSet = "compare_and_set"
	Delete        = "delete"
	AddElement    = "add_element"
)

type ClusterInfo struct {
	CurrentTerm uint32   `validate:"required" json:"currentTerm"`
	CommitIndex uint64   `validate:"required" json:"commitIndex"`
	LastApplied uint64   `validate:"required" json:"lastApplied"`
	NextIndex   []uint64 `validate:"required" json:"nextIndex"`
	MatchIndex  []uint64 `validate:"required" json:"matchIndex"`
}

type LogCommand struct {
	Id       string `validate:"required" json:"id" format:"uuid"`
	Key      string `validate:"required" json:"key"`
	SubKey   string `validate:"required" json:"subKey"`
	OldValue any    `json:"oldValue"`
	NewValue any    `json:"newValue"`
	Type     string `enums:"set,compare_and_set,delete,add_element" validate:"required" json:"type"`
}

type CommandExecutionInfo struct {
	Found   bool   `validate:"required" json:"found"`
	Value   any    `validate:"required" json:"value"`
	Message string `validate:"required" json:"message"`
	Success bool   `validate:"required" json:"success"`
}

type LogEntry struct {
	Term    uint32     `validate:"required" json:"term"`
	Command LogCommand `validate:"required" json:"command"`
}

type ErrorResponse struct {
	Error string `validate:"required" json:"error"`
}

type SetKeyValueRequest struct {
	Value any `json:"value"`
}

type CompareAndSetKeyValueRequest struct {
	OldValue any `json:"oldValue"`
	NewValue any `json:"newValue"`
}

type CommandResponse struct {
	IsLeader  bool   `validate:"required" json:"isLeader"`
	LeaderId  string `validate:"required" json:"leaderId"`
	RequestId string `validate:"required" json:"requestId" format:"uuid"`
}

type GetClusterInfoResponse struct {
	IsLeader bool         `validate:"required" json:"isLeader"`
	LeaderId string       `validate:"required" json:"leaderId"`
	Info     *ClusterInfo `json:"info"`
}

type GetLogResponse struct {
	IsLeader bool       `validate:"required" json:"isLeader"`
	LeaderId string     `validate:"required" json:"leaderId"`
	Entries  []LogEntry `validate:"required" json:"entries"`
}

type GetCommandExecutionInfoResponse struct {
	IsLeader bool                 `validate:"required" json:"isLeader"`
	LeaderId string               `validate:"required" json:"leaderId"`
	Info     CommandExecutionInfo `validate:"required" json:"info"`
}
