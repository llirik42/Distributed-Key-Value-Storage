package common

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/log"
)

type RequestHandler struct {
	ctx *context.Context
}

func NewRequestHandler(ctx *context.Context) *RequestHandler {
	return &RequestHandler{
		ctx: ctx,
	}
}

func (handler *RequestHandler) SetKey(request *SetKeyRequest) (*SetKeyResponse, error) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	leaderId := ctx.GetLeaderId()

	if !ctx.IsLeader() {
		return &SetKeyResponse{
			Code:     NotLeader,
			LeaderId: leaderId,
		}, nil
	}

	cmd := setKeyRequestToCommand(request)
	ctx.PushCommand(cmd)

	return &SetKeyResponse{
		Code:     Success,
		LeaderId: leaderId,
	}, nil
}

func (handler *RequestHandler) GetKey(request *GetKeyRequest) (*GetKeyResponse, error) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	keyValueStorage := ctx.GetKeyValueStorage()
	value := keyValueStorage.Get(request.Key)

	var code string
	if ctx.IsLeader() {
		code = Success
	} else {
		code = NotLeader
	}

	return &GetKeyResponse{
		Value:    value,
		Code:     code,
		LeaderId: ctx.GetLeaderId(),
	}, nil
}

func (handler *RequestHandler) DeleteKey(request *DeleteKeyRequest) (*DeleteKeyResponse, error) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	leaderId := ctx.GetLeaderId()

	if !ctx.IsLeader() {
		return &DeleteKeyResponse{
			Code:     NotLeader,
			LeaderId: leaderId,
		}, nil
	}

	cmd := deleteKeyRequestToCommand(request)
	ctx.PushCommand(cmd)

	return &DeleteKeyResponse{
		Code:     Success,
		LeaderId: leaderId,
	}, nil
}

func (handler *RequestHandler) GetClusterInfo(_ *GetClusterInfoRequest) (*GetClusterInfoResponse, error) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	leaderId := ctx.GetLeaderId()

	if !ctx.IsLeader() {
		return &GetClusterInfoResponse{
			Code:     NotLeader,
			LeaderId: leaderId,
			Info: struct {
				CurrentTerm uint32
				CommitIndex uint64
				LastApplied uint64
				NextIndex   []uint64
				MatchIndex  []uint64
			}{},
		}, nil
	}

	return &GetClusterInfoResponse{
		Code:     Success,
		LeaderId: leaderId,
		Info: struct {
			CurrentTerm uint32
			CommitIndex uint64
			LastApplied uint64
			NextIndex   []uint64
			MatchIndex  []uint64
		}{
			CurrentTerm: ctx.GetCurrentTerm(),
			CommitIndex: ctx.GetCommitIndex(),
			LastApplied: ctx.GetLastApplied(),
			NextIndex:   ctx.GetNextIndexes(),
			MatchIndex:  ctx.GetMatchIndexes(),
		},
	}, nil
}

func (handler *RequestHandler) GetLog(_ *GetLogRequest) (*GetLogResponse, error) {
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	logStorage := ctx.GetLogStorage()
	leaderId := ctx.GetLeaderId()
	entries := logStorage.GetLogEntries(1) // Get all entries

	var code string
	if ctx.IsLeader() {
		code = Success
	} else {
		code = NotLeader
	}

	return &GetLogResponse{
		Code:     code,
		LeaderId: leaderId,
		entries:  entries,
	}, nil
}

func setKeyRequestToCommand(request *SetKeyRequest) log.Command {
	return log.Command{
		Key:   request.Key,
		Value: request.Value,
		Type:  log.Set,
	}
}

func deleteKeyRequestToCommand(request *DeleteKeyRequest) log.Command {
	return log.Command{
		Key:  request.Key,
		Type: log.Delete,
	}
}
