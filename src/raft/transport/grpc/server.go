package grpc

import (
	"context"
	pb "distributed-key-value-storage/generated/proto"
	"distributed-key-value-storage/src/raft/transport"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	pb.RaftServiceServer
	gRPCServer                  *grpc.Server
	listener                    net.Listener
	handleRequestForVoteRequest transport.HandleRequestForVoteRequest
	handleAppendEntriesRequest  transport.HandleAppendEntriesRequest
}

func (server *Server) RequestForVote(
	_ context.Context,
	pbRequest *pb.RequestVoteRequest,
) (*pb.RequestVoteResponse, error) {
	req := MapRequestForVoteRequestFromGRPC(pbRequest)
	resp, err := server.handleRequestForVoteRequest(req)

	if err != nil {
		return nil, fmt.Errorf("failed to handle request for vote %s: %w", pbRequest.String(), err)
	}

	return MapRequestForVoteResponseToGRPC(resp), nil
}

func (server *Server) AppendEntries(
	_ context.Context,
	pbRequest *pb.AppendEntriesRequest,
) (*pb.AppendEntriesResponse, error) {
	req := MapAppendEntriesRequestFromGRPC(pbRequest)
	resp, err := server.handleAppendEntriesRequest(req)

	if err != nil {
		return nil, fmt.Errorf("failed to handle append-entries request %s: %w", pbRequest.String(), err)
	}

	return MapAppendEntriesResponseToGRPC(resp), nil
}

func (server *Server) Listen() error {
	if err := server.gRPCServer.Serve(server.listener); err != nil {
		return fmt.Errorf("failed to serve gRPC-server: %w", err)
	}

	return nil
}

func (server *Server) Shutdown() error {
	server.gRPCServer.GracefulStop()
	return nil
}
