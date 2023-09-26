package grpc_server

import (
	"api-channel/proto"
	"os"
	"time"

	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type (
	Connection struct {
		stream proto.Broadcast_CreateStreamServer
		user   *proto.User
		room   *proto.Room
		active bool
		error  chan error
	}

	FollowerConnection struct {
		stream proto.Broadcast_FollowServer
		server ServerInfo
		active bool
		error  chan error
	}

	Token struct {
		userId     string
		room       string
		expireDate time.Time
	}

	Server struct {
		Id          string
		IsMaster    bool
		Connections map[string]*Connection
		Tokens      map[string]*Token
		Rooms       map[string][]string
		Followers   map[string]*FollowerConnection
		Fserver     *FollowedServer
	}

	ServerInfo struct {
		Id string
	}
)

var (
	server   *Server
	serverId string
)

type GRPCServer struct{}

func (g *GRPCServer) NewChatServer(sid string) *Server {
	serverId = sid
	server = &Server{
		Id:          serverId,
		IsMaster:    false,
		Connections: make(map[string]*Connection),
		Tokens:      make(map[string]*Token),
		Rooms:       make(map[string][]string),
		Followers:   make(map[string]*FollowerConnection),
	}
	return server
}
