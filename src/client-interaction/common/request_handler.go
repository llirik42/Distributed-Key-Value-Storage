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

func (handler *RequestHandler) HandleSetKey(request *SetKeyRequest) (*SetKeyResponse, error) {
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

func (handler *RequestHandler) HandleGetKey(request *GetKeyRequest) (*GetKeyResponse, error) {
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

func (handler *RequestHandler) HandleDeleteKey(request *DeleteKeyRequest) (*DeleteKeyResponse, error) {
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
