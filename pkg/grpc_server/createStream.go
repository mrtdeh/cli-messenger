package grpc_server

import (
	"api-channel/pkg/data"
	"api-channel/proto"
	"fmt"
	"log"
	"time"
)

func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	// prepare new connection info
	conn := &Connection{
		stream: stream,
		user:   pconn.User,
		room:   pconn.Room,
		active: true,
		error:  make(chan error),
	}

	username := pconn.User.Username
	userId := pconn.User.Id
	room := pconn.Room.Name
	// deny requested username if it assigned and actived by another client
	uExist, uActive := s.checkUserExist(username)
	if c, ok := s.Connections[userId]; ok || uExist {
		err := fmt.Errorf("this username(%s) has already been used", pconn.User.Username)
		if (c != nil && c.active) || uActive {
			stream.Send(errorResponse(err))
			return err
		}
	}

	s.updateRooms(string(room))
	s.updateUsers(username, true)
	// add user to connection list
	s.Connections[userId] = conn
	// generate token for user
	newtoken := genToken()
	// add token to token list with user id and expiration date
	s.Tokens[newtoken] = &Token{
		userId:     userId,
		room:       room,
		expireDate: time.Now().Add(time.Second * 60 * 60), // 1h to expire
	}
	// send token for user
	err := stream.Send(tokenResponse(newtoken))
	if err != nil {
		log.Println("error in send token : ", err.Error())
		conn.error <- err
	}
	// s.Rooms = map[RoomName][]UserId{}
	s.addUserToRoom(userId, room)
	// send join message for all users in lobby
	s.sendToServer(string(room), "", "", userStatusResponse(pconn.User.Username, proto.UserStatus_join))

	return <-conn.error
}

func (s *Server) checkUserExist(username string) (exist bool, active bool) {
	if s.IsMaster {
		if u, ok := data.Users[username]; ok {
			exist = true
			active = u.Active
		}
	} else {
		exist, active = s.Fserver.userExist(username)
	}

	return exist, active
}

func (s *Server) addUserToRoom(userId string, room string) {
	if users, ok := s.Rooms[room]; ok {
		for _, u := range users {
			if u == userId {
				return
			}
		}
		users = append(users, userId)
		s.Rooms[room] = users
	} else {
		s.Rooms[room] = []string{userId}
	}
}
