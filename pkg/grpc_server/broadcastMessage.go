package grpc_server

import (
	"api-channel/proto"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func RemoveString(slice []string, str string) []string {
	var result []string
	for _, s := range slice {
		if s != str {
			result = append(result, s)
		}
	}
	return result
}

func (s *Server) BroadcastMessage(stream proto.Broadcast_BroadcastMessageServer) error {

	var fInfo *proto.FileInfo
	var file *os.File
	var token string
	var room string

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			if t := s.Tokens[token]; t != nil {
				u := t.userId
				s.Connections[u].active = false
				username := s.Connections[u].user.Username
				s.updateUsers(username, false)

				user := s.Connections[u].user.Username
				room := s.Connections[u].room.Name
				s.Rooms[room] = RemoveString(s.Rooms[room], u)
				s.sendToServer(string(room), "", "", userStatusResponse(user, proto.UserStatus_leave))
			}
			log.Printf("Error while reading client stream: %v\n", err)
			return err
		}

		token = req.Token
		msg := req.GetMsg()

		room = "lobby"
		if t, ok := s.Tokens[token]; ok {
			room = t.room
		}

		if t := msg.GetTextMsg(); t != nil {
			// fmt.Println("debug get message")
			res := s.textMessage(req, t.Content)
			s.sendToServer(string(room), "", "", res)

		} else if f := msg.GetFileMsg(); f != nil {

			if info := f.GetInfo(); info != nil {
				name := path.Base(info.Name)
				file, err = os.Create("/tmp/" + name)
				if err != nil {
					log.Println("error in open file : ", err.Error())
					return err
				}
			} else if chunk := f.GetChunkData(); chunk != nil {
				_, err := file.Write(chunk)
				if err != nil {
					log.Println("error in write to file : ", err.Error())
					return err
				}
			} else if done := f.GetDone(); done {
				file.Close()
			} else {
				log.Fatal("Oops! command not found")
			}

			// store file info
			if info := f.GetInfo(); info != nil {
				fInfo = info
			}
			// send file info to all connections
			if done := f.GetDone(); done {
				res := s.textMessage(req, fmt.Sprintf("File[name=%s, size=%d]", fInfo.Name, fInfo.Size))
				s.sendToServer(string(room), "", "", res)
				break
			}
		}

	}

	return nil
}
