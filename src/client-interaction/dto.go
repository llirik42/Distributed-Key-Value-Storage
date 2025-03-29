package client_interaction

const (
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

type LogEntry struct {
	Term    uint32     `validate:"required" json:"term"`
	Command LogCommand `validate:"required" json:"command"`
}

type ErrorResponse struct {
	Error string `validate:"required" json:"error"`
}

type SetKeyRequest struct {
	Value any `validate:"required" json:"value"`
}

type SetKeyResponse struct {
	IsLeader  bool   `validate:"required" json:"isLeader"`
	LeaderId  string `validate:"required" json:"leaderId"`
	RequestId string `validate:"required" json:"requestId" format:"uuid"`
}

type GetKeyResponse struct {
	IsLeader bool   `validate:"required" json:"isLeader"`
	Value    any    `validate:"required" json:"value"`
	Code     string `enums:"success,not_found" validate:"required" json:"code"`
	LeaderId string `validate:"required" json:"leaderId"`
}

type DeleteKeyResponse struct {
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
