package grpc_client

import (
	"api-channel/proto"
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
)

func SendFile(filepath string) error {
	if myToken == "" {
		return fmt.Errorf("token not found in client")
	}

	var err error
	var frec *os.File
	var stat fs.FileInfo
	// open selected file
	frec, err = os.Open(filepath)
	if err != nil {
		log.Fatal("cannot open txt file: ", err)
	}
	defer frec.Close()

	stat, err = frec.Stat()
	if err != nil {
		log.Fatal("file read failed : ", err)
	}
	// failed when is directory
	if stat.IsDir() {
		return fmt.Errorf("you select a dir, but must select a file")
	}

	reader := bufio.NewReader(frec)
	buffer := make([]byte, 1024)

	// connnect to RPC function
	s, err := client.BroadcastMessage(context.Background())
	if err != nil {
		fmt.Printf("http api -> Error Sending Info: %v", err)
	}

	info := &proto.MessageRequest{
		Token: myToken,
	}

	info.Msg = &proto.Message{
		Data: &proto.Message_FileMsg{
			FileMsg: &proto.FileMessage{
				Data: &proto.FileMessage_Info{
					Info: &proto.FileInfo{
						Name: frec.Name(),
						Size: uint32(stat.Size()),
					},
				},
			},
		},
	}

	// send file info to server
	s.Send(info)

	// send file chunk data to server
	for {
		// read file buffer
		n, err := reader.Read(buffer)
		// check if end of file
		if err == io.EOF {
			break
		}
		// check error
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		info.Msg = &proto.Message{
			Data: &proto.Message_FileMsg{
				FileMsg: &proto.FileMessage{
					Data: &proto.FileMessage_ChunkData{
						ChunkData: buffer[:n],
					},
				},
			},
		}
		// send chunk
		s.Send(info)
	}

	info.Msg = &proto.Message{
		Data: &proto.Message_FileMsg{
			FileMsg: &proto.FileMessage{
				Data: &proto.FileMessage_Done{
					Done: true,
				},
			},
		},
	}
	//
	// send done flag to server
	s.Send(info)

	return nil
}
