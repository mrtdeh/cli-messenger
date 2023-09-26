package grpc_server

import (
	"api-channel/pkg/data"
	"api-channel/proto"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
)

func (s *Server) forward(srv, room string, msg *proto.ServerResponse) {
	if f, ok := s.Followers[srv]; ok {
		if f.active {
			f.stream.Send(&proto.LeaderResponse{
				Data: &proto.LeaderResponse_ForwardResponse{
					ForwardResponse: &proto.ForwardResponse{
						ServerRes: msg,
						Room:      room,
					},
				},
			})
		} else {
			fmt.Printf("follower id=%s not active\n", srv)
		}
	} else {
		fmt.Printf("follower id=%s not found\n", srv)
	}
}

func (s *Server) updateUsers(username string, active bool) {
	if s.IsMaster {
		data.AddUser(username, active, s.Id)
	} else {
		users := []*proto.User{
			{
				Username: username,
				Active:   active,
			},
		}
		s.Fserver.send(&proto.FollowerRequest{
			Data: &proto.FollowerRequest_UpdateMsg{
				UpdateMsg: &proto.UpdateMessage{
					Data: &proto.UpdateMessage_NewUsers{
						NewUsers: &proto.NewUsers{
							Users: users,
						},
					},
				},
			},
		})
	}
}

func (s *Server) updateRooms(room string) {
	r := string(room)

	data.AddServerToRoom(r, serverId)
	if s.IsMaster {
		s.sendToFollowers(&proto.LeaderResponse{
			Data: &proto.LeaderResponse_UpdateMsg{
				UpdateMsg: &proto.UpdateMessage{
					Data: &proto.UpdateMessage_NewRooms{
						NewRooms: &proto.NewRooms{
							Rooms: []*proto.Room{
								{
									Name:     string(room),
									ServerId: serverId,
								},
							},
						},
					},
				},
			},
		})
	} else {
		fmt.Printf("debug send server id to master room=%s , serverId=%s\n", room, serverId)
		s.Fserver.send(&proto.FollowerRequest{
			Data: &proto.FollowerRequest_UpdateMsg{
				UpdateMsg: &proto.UpdateMessage{
					Data: &proto.UpdateMessage_NewRooms{
						NewRooms: &proto.NewRooms{
							Rooms: []*proto.Room{
								{
									Name:     string(room),
									ServerId: serverId,
								},
							},
						},
					},
				},
			},
		})
	}
}

func (s *Server) sendToAll(msg *proto.ServerResponse) {
	wait := sync.WaitGroup{}
	done := make(chan int)
	grpcLog.Info("Sending message to: all")
	for _, conn := range s.Connections {
		wait.Add(1)

		go func(req *proto.ServerResponse, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(req)

				if err != nil {
					grpcLog.Errorf("Error with Stream: %v - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}

		}(msg, conn)

	}
	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}

func (s *Server) sendToFollowers(msg *proto.LeaderResponse) {
	wait := sync.WaitGroup{}
	done := make(chan int)
	grpcLog.Info("Sending message to followers")

	for _, v := range s.Followers {
		wait.Add(1)

		go func(req *proto.LeaderResponse, fconn *FollowerConnection) {
			defer wait.Done()

			if fconn.active {
				err := fconn.stream.Send(req)

				if err != nil {
					grpcLog.Errorf("Error with Stream: %v - Error: %v", fconn.stream, err)
					fconn.active = false
					fconn.error <- err
				}
			}

		}(msg, v)

	}
	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}

func (s *Server) sendToServer(room, from, to string, msg *proto.ServerResponse) {

	var servers []string
	var ok bool

	if servers, ok = data.RoomsServers[room]; !ok {
		log.Println("room not found : ", room)
		return
	}

	// iterate servers of room
	for _, srv := range servers {
		// if srv not current server
		if srv != serverId {
			// ignore forward if srv is requested server
			if srv == from {
				continue
			}

			if to != "" && to != srv {
				continue
			}
			// if me is master
			if s.IsMaster {
				// forward to other (followers)
				s.forward(srv, room, msg)
			} else {
				// else, forward to master
				s.Fserver.send(&proto.FollowerRequest{
					Data: &proto.FollowerRequest_ForwardMsg{
						ForwardMsg: &proto.ForwardRequest{
							From:      serverId,
							To:        srv,
							ServerRes: msg,
							Room:      room,
						},
					},
				})
			}

		} else {
			if s.IsMaster && to != "" && to != serverId {
				continue
			}
			s.sendToRoom(room, msg)
		}
	}
}

func (s *Server) textMessage(req *proto.MessageRequest, msg string) *proto.ServerResponse {
	var from *proto.User
	if t, ok := s.Tokens[req.Token]; ok {
		if u, ok := s.Connections[t.userId]; ok {
			from = u.user
		}
	}
	return msgResponse(msg, from)
}

func errorResponse(err error) *proto.ServerResponse {
	return &proto.ServerResponse{
		Data: &proto.ServerResponse_ErrorResponse{
			ErrorResponse: &proto.ErrorResponse{
				Error: err.Error(),
			},
		},
	}
}

func msgResponse(msg string, from *proto.User) *proto.ServerResponse {
	return &proto.ServerResponse{
		Data: &proto.ServerResponse_MsgResponse{
			MsgResponse: &proto.MessageResponse{
				From:      from,
				Timestamp: time.Now().String(),
				Msg: &proto.Message{
					Data: &proto.Message_TextMsg{
						TextMsg: &proto.TextMessage{
							Content: msg,
						},
					},
				},
			},
		},
	}
}

func tokenResponse(token string) *proto.ServerResponse {
	return &proto.ServerResponse{
		Data: &proto.ServerResponse_TokenResponse{
			TokenResponse: &proto.TokenResponse{
				Token: token,
			},
		},
	}
}

func userStatusResponse(user string, s proto.UserStatus) *proto.ServerResponse {
	return &proto.ServerResponse{
		Data: &proto.ServerResponse_UserStatusResponse{
			UserStatusResponse: &proto.UserStatusResponse{
				Name:       user,
				UserStatus: s,
			},
		},
	}
}

func (s *Server) sendToRoom(room string, msg *proto.ServerResponse) {
	wait := sync.WaitGroup{}
	done := make(chan int)
	grpcLog.Info("Sending message to: ", room)

	if _, ok := s.Rooms[room]; !ok {
		log.Println("room not found : ", room)
		return
	}

	for _, uid := range s.Rooms[room] {
		wait.Add(1)

		conn := s.Connections[uid]

		go func(req *proto.ServerResponse, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(req)

				if err != nil {
					grpcLog.Errorf("Error with Stream: %v - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}

		}(msg, conn)

	}
	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}

func genToken() string {
	md5Sum := md5.Sum([]byte(time.Now().String()))
	return hex.EncodeToString(md5Sum[:])
}
