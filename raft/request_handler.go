package raft

import (
	"distributed-algorithms/raft/context"
	"distributed-algorithms/raft/dto"
)

type RequestHandler struct {
	ctx *node.Context
}

func NewRequestHandler(ctx *node.Context) *RequestHandler {
	return &RequestHandler{ctx: ctx}
}

func (handler *RequestHandler) HandleRequestVoteRequest(request dto.RequestVoteRequest) (*dto.RequestVoteResponse, error) {
	ctx := handler.ctx
	ctx.CheckTerm(request.Term)

}

func (handler *RequestHandler) HandleAppendEntriesRequest(request dto.AppendEntriesRequest) (*dto.AppendEntriesResponse, error) {
	ctx := handler.ctx
	ctx.CheckTerm(request.Term)
}
