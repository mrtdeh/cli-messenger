package grpc_client

import (
	"api-channel/proto"
	"fmt"
	"time"
)

var clientStream proto.Broadcast_BroadcastMessageClient

func SendText(text string) error {
	if myToken == "" {
		return fmt.Errorf("token not found in client")
	}

	// var err error

	msg := &proto.MessageRequest{
		Token:     myToken,
		Timestamp: time.Now().String(),
		Msg: &proto.Message{
			Data: &proto.Message_TextMsg{
				TextMsg: &proto.TextMessage{
					Content: text,
				},
			},
		},
	}

	// if clientStream == nil {
	// 	var err error
	// 	clientStream, err = client.BroadcastMessage(context.Background())
	// 	if err != nil {
	// 		return fmt.Errorf("error in create client streaming for send text : %s", err.Error())
	// 	}
	// }

	err := clientStream.Send(msg)
	if err != nil {
		return fmt.Errorf("error in send text : %s", err.Error())
	}

	return nil

}
