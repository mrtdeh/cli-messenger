package main

import (
	"api-channel/pkg/cli"
	grpc_client "api-channel/pkg/grpc_client"
	"api-channel/proto"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"encoding/hex"
	"log"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	roomName *string
)

func main() {

	srvHost := flag.String("shost", "localhost", "server host")
	srvPort := flag.String("sport", "8082", "server port")
	name := flag.String("n", "mory", "The name of the user")
	roomName = flag.String("r", "lobby", "The name of the room")
	flag.Parse()

	hostPort := *srvHost + ":" + *srvPort
	conn, err := grpc.Dial(hostPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldnt connect to service: %v", err)
	}

	grpc_client.SetClient(proto.NewBroadcastClient(conn))

	id := sha256.Sum256([]byte(*name))

	// s1 := rand.NewSource(time.Now().UnixNano())
	// r1 := rand.New(s1)
	colors := []string{"white", "red", "green", "magenta", "cyan"}
	user := &proto.User{
		Id:       hex.EncodeToString(id[:]),
		Username: *name,
		Color:    colors[0], //colors[r1.Intn(len(colors))],
	}
	room := &proto.Room{
		Name: *roomName,
	}
	err = grpc_client.Connect(user, room)
	if err != nil {
		log.Fatal(err)
	}

	err = grpc_client.WaitForOK()
	if err != nil {
		log.Fatal(err)
	}

	go serverHealthCheck()

	// ========================= Liner configuration ==============================
	Prompt(user.Username, user.Color)

}

func serverHealthCheck() {
	for {
		err := grpc_client.Ping()
		if err != nil {
			fmt.Println()
			fmt.Println("Ping error : ", err.Error())
			cli.Close()
		}
		time.Sleep(time.Second)
	}
}

func Prompt(user, c string) {
	var err error

	cc := cli.GetColorCode(c)
	col := color.New(cc).SprintFunc()

	prompt := fmt.Sprintf("%s@%s ", *roomName, col(user))
	l, err := readline.NewEx(&readline.Config{
		Prompt:          prompt,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		Listener:          readline.FuncListener(cli.ListernerReadline),
		HistorySearchFold: true,
		// FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}

	cli.SetReadLine(l)
	cli.SetRoome(*roomName)

	l.CaptureExitSignal()
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		promptResponse(line)
	}
}

func promptResponse(str string) {
	str = strings.Trim(str, " ")
	if str == "" {
		return
	}
	action := str[:1]

	if action == "@" {
		// Response to send file
		err := grpc_client.SendFile(str[1:])
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		// Response to send text
		err := grpc_client.SendText(str)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
