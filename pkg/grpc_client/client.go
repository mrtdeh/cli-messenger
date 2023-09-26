package grpc_client

import (
	"api-channel/pkg/cli"
	"api-channel/proto"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	client proto.BroadcastClient
	wait   *sync.WaitGroup

	user *proto.User

	myToken string
)

func init() {
	wait = &sync.WaitGroup{}
}

func SetClient(c proto.BroadcastClient) {
	client = c
}

func Ping() error {
	res, err := client.Ping(context.Background(), &proto.PingRequest{})
	if err != nil {
		return err
	}
	if res == nil {
		return fmt.Errorf("server not available")
	}

	return nil
}

func WaitForOK() error {
	timeout := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-timeout.C:
			return fmt.Errorf("connect to server timeout")
		default:
			if myToken != "" {
				return nil
			}
		}
	}
}

func Connect(user *proto.User, room *proto.Room) error {
	var streamerror error
	// user = u
	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User:   user,
		Room:   room,
		Active: true,
	})

	if err != nil {
		fmt.Printf("connection failed: %v\n", err)
		os.Exit(0)
	}

	clientStream, err = client.BroadcastMessage(context.Background())
	if err != nil {
		return fmt.Errorf("error in create client streaming for send text : %s", err.Error())
	}

	wait.Add(1)
	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()
		var file *os.File

		for {
			req, err := str.Recv()
			if err == io.EOF {
				file.Close()
				continue
			}
			if err != nil {
				streamerror = fmt.Errorf("error reading message: %w", err)
				break
			}

			// error response
			if e := req.GetErrorResponse(); e != nil {
				fmt.Println(e.Error)
				os.Exit(0)
			}

			// token response
			if t := req.GetTokenResponse(); t != nil {
				myToken = t.Token
				clientStream.Send(&proto.MessageRequest{Token: myToken})
			}

			// join response
			if s := req.GetUserStatusResponse(); s != nil {
				if s.Name != user.Username {
					if s.UserStatus == proto.UserStatus_join {
						green := color.New(color.FgHiGreen).SprintFunc()

						cli.Print(green(s.Name + " joined!"))
					} else if s.UserStatus == proto.UserStatus_leave {
						red := color.New(color.FgHiRed).SprintFunc()
						cli.Print(red(s.Name + " leaved!"))
					}
				}
			}

			// msg response
			if msg := req.GetMsgResponse(); msg != nil {
				if m := msg.Msg.GetTextMsg(); m != nil {
					username := msg.From.Username
					color := msg.From.Color
					if username != user.Username {
						cli.PrintUserMessage(username, color, m.Content)
					}
				}

			}
		}
	}(stream)

	return streamerror
}
