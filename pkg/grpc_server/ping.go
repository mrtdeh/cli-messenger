package grpc_server

import (
	"api-channel/proto"
	"context"
)

func (s *Server) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {

	return &proto.PingResponse{}, nil
}
