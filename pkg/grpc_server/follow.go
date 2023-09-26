package grpc_server

import (
	"api-channel/pkg/data"
	"api-channel/proto"
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FollowedServer struct {
	Id       string
	IsMaster bool
	conn     proto.BroadcastClient
	req      proto.Broadcast_FollowClient
}

func (s *Server) FollowMaster(fs *FollowedServer) error {
	s.Fserver = fs

	req, err := fs.conn.Follow(context.Background())
	if err != nil {
		return fmt.Errorf("failed to follow : %s\n", err.Error())
	}
	fs.req = req

	req.Send(&proto.FollowerRequest{
		Data: &proto.FollowerRequest_JoinMsg{
			JoinMsg: &proto.JoinMessage{
				Id: serverId,
			},
		},
	})

	for {
		res, err := req.Recv()
		if err != nil {
			return fmt.Errorf("error in recivce : %s", err.Error())
		}

		if u := res.GetUpdateMsg(); u != nil {
			if r := u.GetNewRooms(); r != nil {
				for _, rr := range r.Rooms {
					data.AddServerToRoom(rr.Name, rr.ServerId)
				}
			}
		} else if f := res.GetForwardResponse(); f != nil {
			// fmt.Println("forwarded from master")
			server.sendToRoom(f.Room, f.ServerRes)
		}

	}
}

// Check user exit in master server
//
//	note: master server already store users data from servers in the cluster
func (fs *FollowedServer) userExist(username string) (bool, bool) {
	var exist, active bool

	e, _ := fs.conn.UserExist(context.Background(), &proto.UserExistRequest{
		Username: username,
	})
	exist = e.Exist
	active = e.Active

	return exist, active
}

// Send follower request to master server
func (fs *FollowedServer) send(msg *proto.FollowerRequest) error {
	err := fs.req.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

type GRPCClient struct{}

// conenct to specified address if this address is for master server
func (g *GRPCClient) ConnectIfMaster(addr string) (*FollowedServer, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error in dial : %s", err.Error())
	}
	c := proto.NewBroadcastClient(conn)
	im, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
	if err != nil {
		return nil, fmt.Errorf("error check is_master : %s\n", err.Error())
	}

	// fmt.Printf("the server %s is master : %t\n", addr, im.IsMaster)

	fs := &FollowedServer{
		IsMaster: im.IsMaster,
		Id:       im.Id,
		conn:     c,
	}

	if im.IsMaster == false {
		fmt.Println("close connection from ", addr)
		conn.Close()

		return fs, nil
	}

	return fs, nil
}

func (s *Server) Follow(stream proto.Broadcast_FollowServer) error {

	// forwarded server id
	var fsId string

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("follower connection closed : ", err.Error())
			data.CleanServer(fsId)
			return err
		}

		if j := req.GetJoinMsg(); j != nil {
			fsId = j.Id

			// prepare new connection info
			fconn := &FollowerConnection{
				stream: stream,
				server: ServerInfo{
					Id: fsId,
				},
				active: true,
				error:  make(chan error),
			}
			s.Followers[fsId] = fconn

			fmt.Println("new follower added : ", fsId)

		} else if u := req.GetUpdateMsg(); u != nil {
			if rooms := u.GetNewRooms(); rooms != nil {

				for _, r := range rooms.Rooms {
					data.AddServerToRoom(r.Name, r.ServerId)
				}

				s.sendToFollowers(&proto.LeaderResponse{
					Data: &proto.LeaderResponse_UpdateMsg{
						UpdateMsg: u,
					},
				})

			} else if users := u.GetNewUsers(); users != nil {
				for _, u := range users.Users {
					data.AddUser(u.Username, u.Active, fsId)
				}
			}
		} else if f := req.GetForwardMsg(); f != nil {
			s.sendToServer(f.Room, f.From, f.To, f.ServerRes)
		}

	}

	return nil
}
