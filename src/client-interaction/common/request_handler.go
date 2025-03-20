package common

import (
	"distributed-algorithms/src/context"
	"distributed-algorithms/src/key-value"
)

type RequestHandler struct {
	context         *context.Context
	keyValueStorage *key_value.Storage
}

func NewRequestHandler(context *context.Context, keyValueStorage *key_value.Storage) *RequestHandler {
	return &RequestHandler{
		context:         context,
		keyValueStorage: keyValueStorage,
	}
}

func (handler *RequestHandler) HandleSetKey(request *SetKeyRequest) (*SetKeyResponse, error) {
	// TODO: implement

	return &SetKeyResponse{
		Code:     Success,
		LeaderId: "Node-1",
	}, nil
}

func (handler *RequestHandler) HandleGetKey(request *GetKeyRequest) (*GetKeyResponse, error) {
	// TODO: implement

	return &GetKeyResponse{
		Value:    nil,
		Code:     Success,
		LeaderId: "Node-1",
	}, nil
}

func (handler *RequestHandler) HandleDeleteKey(request *DeleteKeyRequest) (*DeleteKeyResponse, error) {
	// TODO: implement

	return &DeleteKeyResponse{
		Code:     Success,
		LeaderId: "Node-1",
	}, nil
}
