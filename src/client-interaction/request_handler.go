package client_interaction

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/log"
	"distributed-algorithms/src/log/executor"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestHandler struct {
	ctx *context.Context
}

func NewRequestHandler(ctx *context.Context) *RequestHandler {
	return &RequestHandler{
		ctx: ctx,
	}
}

// SetKeyValue
// @Id "SetKeyValue"
// @Router /key/{key} [post]
// @Summary Set Key Value
// @Description Sets value for the given key. If the old value already exists, it is replaced by a new one.
// @Tags key
// @Param key path string true " "
// @Param request body SetKeyValueRequest	true " "
// @Success 200 {object} CommandResponse
// @Failure 400 {object} ErrorResponse
func (handler *RequestHandler) SetKeyValue(c *gin.Context) {
	key := c.Param("key")
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	request := SetKeyValueRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var requestId = ""
	isLeader := ctx.IsLeader()
	if isLeader {
		cmd := createSetKeyCommand(key, &request)
		requestId = ctx.PushCommand(cmd)
	}

	response := CommandResponse{
		IsLeader:  isLeader,
		LeaderId:  ctx.GetLeaderId(),
		RequestId: requestId,
	}

	c.JSON(http.StatusOK, response)
}

// CompareAndSetKeyValue
// @Id "CompareAndSetKeyValue"
// @Router /key/{key} [patch]
// @Summary Compare And Set Key Value
// @Tags key
// @Param key path string true " "
// @Param request body CompareAndSetKeyValueRequest true " "
// @Success 200 {object} CommandResponse
// @Failure 400 {object} ErrorResponse
func (handler *RequestHandler) CompareAndSetKeyValue(c *gin.Context) {
	key := c.Param("key")
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	request := CompareAndSetKeyValueRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var requestId = ""
	isLeader := ctx.IsLeader()
	if isLeader {
		cmd := createCASCommand(key, &request)
		requestId = ctx.PushCommand(cmd)
	}

	response := CommandResponse{
		IsLeader:  isLeader,
		LeaderId:  ctx.GetLeaderId(),
		RequestId: requestId,
	}

	c.JSON(http.StatusOK, response)
}

// GetKeyValue
// @Id "GetKeyValue"
// @Router /key/{key} [get]
// @Summary Get Key Value
// @Tags key
// @Param key path string true " "
// @Success 200 {object} CommandResponse
func (handler *RequestHandler) GetKeyValue(c *gin.Context) {
	key := c.Param("key")
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	var requestId = ""
	isLeader := ctx.IsLeader()
	if isLeader {
		cmd := createGetKeyCommand(key)
		requestId = ctx.PushCommand(cmd)
	}

	response := CommandResponse{
		IsLeader:  isLeader,
		LeaderId:  ctx.GetLeaderId(),
		RequestId: requestId,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteKey
// @Id "DeleteKeyValue"
// @Router /key/{key} [delete]
// @Summary Delete Key Value
// @Description Deletes value for the given key
// @Tags key
// @Param key path string true " "
// @Success 200 {object} CommandResponse
func (handler *RequestHandler) DeleteKey(c *gin.Context) {
	key := c.Param("key")
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	var requestId = ""
	isLeader := ctx.IsLeader()
	if isLeader {
		cmd := createDeleteKeyCommand(key)
		requestId = ctx.PushCommand(cmd)
	}

	response := CommandResponse{
		IsLeader:  isLeader,
		LeaderId:  ctx.GetLeaderId(),
		RequestId: requestId,
	}

	c.JSON(http.StatusOK, response)
}

// GetClusterInfo
// @Id "GetClusterInfo"
// @Router /cluster/info [get]
// @Summary Get Cluster Info
// @Tags cluster
// @Success 200 {object} GetClusterInfoResponse
func (handler *RequestHandler) GetClusterInfo(c *gin.Context) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	var info *ClusterInfo = nil
	isLeader := ctx.IsLeader()

	if isLeader {
		info = &ClusterInfo{
			CurrentTerm: ctx.GetCurrentTerm(),
			CommitIndex: ctx.GetCommitIndex(),
			LastApplied: ctx.GetLastApplied(),
			NextIndex:   ctx.GetNextIndexes(),
			MatchIndex:  ctx.GetMatchIndexes(),
		}
	}

	response := GetClusterInfoResponse{
		IsLeader: isLeader,
		LeaderId: ctx.GetLeaderId(),
		Info:     info,
	}

	c.JSON(http.StatusOK, response)
}

// GetLog
// @Id "GetLog"
// @Router /cluster/log [get]
// @Summary Get Log
// @Tags cluster
// @Success 200 {object} GetLogResponse
func (handler *RequestHandler) GetLog(c *gin.Context) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	logStorage := ctx.GetLogStorage()
	entries := logStorage.GetLogEntries(1) // Get all entries

	response := GetLogResponse{
		IsLeader: ctx.IsLeader(),
		LeaderId: ctx.GetLeaderId(),
		Entries:  mapLogEntries(entries),
	}

	c.JSON(http.StatusOK, response)
}

// GetCommandExecutionInfo
// @Id "GetCommandExecutionInfo"
// @Router /command/{commandId} [get]
// @Summary Get Command Execution Info
// @Tags storage
// @Success 200 {object} GetCommandExecutionInfoResponse
func (handler *RequestHandler) GetCommandExecutionInfo(c *gin.Context) {
	commandId := c.Param("commandId")
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	info, exists := executor.GetCommandExecutionInfo(
		ctx.GetKeyValueStorage(),
		ctx.GetExecutedCommandsKey(),
		commandId,
	)

	response := GetCommandExecutionInfoResponse{
		Found: exists,
		Info: CommandExecutionInfo{
			Value:   info.Value,
			Message: info.Message,
			Success: info.Success,
		},
	}

	c.JSON(http.StatusOK, response)
}

func createGetKeyCommand(key string) log.Command {
	return log.Command{
		Key:  key,
		Type: log.Get,
	}
}

func createSetKeyCommand(key string, request *SetKeyValueRequest) log.Command {
	return log.Command{
		Key:      key,
		NewValue: request.Value,
		Type:     log.Set,
	}
}

func createCASCommand(key string, request *CompareAndSetKeyValueRequest) log.Command {
	return log.Command{
		Key:      key,
		OldValue: request.OldValue,
		NewValue: request.NewValue,
		Type:     log.CompareAndSet,
	}
}

func createDeleteKeyCommand(key string) log.Command {
	return log.Command{
		Key:  key,
		Type: log.Delete,
	}
}

func mapLogEntries(entries []log.Entry) []LogEntry {
	result := make([]LogEntry, len(entries))

	for i, v := range entries {
		cmd := v.Command

		result[i] = LogEntry{
			Term: v.Term,
			Command: LogCommand{
				Id:       cmd.Id,
				Key:      cmd.Key,
				SubKey:   cmd.SubKey,
				OldValue: cmd.OldValue,
				NewValue: cmd.NewValue,
				Type:     mapCommandType(cmd.Type),
			},
		}
	}

	return result
}

func mapCommandType(cmdType int) string {
	switch cmdType {
	case log.Get:
		return Get
	case log.Set:
		return Set
	case log.CompareAndSet:
		return CompareAndSet
	case log.Delete:
		return Delete
	case log.AddElement:
		return AddElement
	default:
		panic(fmt.Errorf("unknown command type: %d", cmdType))
	}
}
