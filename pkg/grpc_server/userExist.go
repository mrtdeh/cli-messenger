package grpc_server

import (
	"api-channel/pkg/data"
	"api-channel/proto"
	"context"
)

func (s *Server) UserExist(ctx context.Context, req *proto.UserExistRequest) (*proto.UserExistResponse, error) {
	var exist, active bool
	username := req.Username

	if u, ok := data.Users[username]; ok {
		exist = true
		active = u.Active
	}

	res := &proto.UserExistResponse{
		Exist:  exist,
		Active: active,
	}
	return res, nil
}
