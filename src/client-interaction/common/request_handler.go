package common

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/key-value"
	"distributed-algorithms/src/log"
	"fmt"
)

type RequestHandler struct {
	ctx     *context.Context
	storage key_value.Storage
}

func NewRequestHandler(context *context.Context, keyValueStorage *key_value.Storage) *RequestHandler {
	return &RequestHandler{
		ctx:     context,
		storage: *keyValueStorage,
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

	if err := ctx.PushCommand(&cmd); err != nil {
		return nil, fmt.Errorf("failed to push SET-KEY command to log: %v", err)
	}

	return &SetKeyResponse{
		Code:     Success,
		LeaderId: leaderId,
	}, nil
}

func (handler *RequestHandler) GetKey(request *GetKeyRequest) (*GetKeyResponse, error) {
	storage := handler.storage
	ctx := handler.ctx
	ctx.Lock()
	defer ctx.Unlock()

	value, err := storage.Get(request.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve value of key: %v", err)
	}

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

	if err := ctx.PushCommand(&cmd); err != nil {
		return nil, fmt.Errorf("failed to push DELETE-KEY command to log: %v", err)
	}

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
			Info:     nil,
		}, nil
	}

	return &GetClusterInfoResponse{
		Code:     Success,
		LeaderId: leaderId,
		Info: &struct {
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

	leaderId := ctx.GetLeaderId()

	if !ctx.IsLeader() {
		return &GetLogResponse{
			Code:     NotLeader,
			LeaderId: leaderId,
			entries:  nil,
		}, nil
	}

	entries, err := ctx.GetLog().GetLogEntries(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get log-entries: %v", err)
	}

	return &GetLogResponse{
		Code:     Success,
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
