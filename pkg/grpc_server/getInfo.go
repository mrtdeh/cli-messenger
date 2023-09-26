package grpc_server

import (
	"api-channel/proto"
	"context"
)

func (s *Server) GetInfo(context.Context, *proto.EmptyRequest) (*proto.InfoResponse, error) {
	return &proto.InfoResponse{
		Id:       s.Id,
		IsMaster: s.IsMaster,
	}, nil
}
