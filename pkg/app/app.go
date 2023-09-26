package app

import (
	"api-channel/pkg/grpc_server"
	"api-channel/proto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	"google.golang.org/grpc"
)

type ServerInterface interface {
	NewChatServer(serverId string) *grpc_server.Server
}
type ClientInterface interface {
	ConnectIfMaster(addr string) (*grpc_server.FollowedServer, error)
}

type AppOptions struct {
	MaxTimeout           time.Duration
	MaxConnectTry        int
	MaxConnectTryTimeout time.Duration
}
type App struct {
	Id            string
	Addr          string
	ReplicasAddrs []string
	Server        ServerInterface
	Client        ClientInterface
	cs            *grpc_server.Server
	gs            *grpc.Server

	status string

	Options *AppOptions
}

func (a *App) setDefaultOptions() {
	if a.Options == nil {
		a.Options = &AppOptions{}
	}
	if a.Options.MaxConnectTry == 0 {
		a.Options.MaxConnectTry = 10
	}
	if a.Options.MaxTimeout == 0 {
		a.Options.MaxTimeout = time.Second * 5
	}
	if a.Options.MaxConnectTryTimeout == 0 {
		a.Options.MaxConnectTryTimeout = time.Second * 1
	}
}

func (a *App) addServer(ad string) {
	for _, addr := range a.ReplicasAddrs {
		if addr == ad {
			return
		}
	}
	a.ReplicasAddrs = append(a.ReplicasAddrs, ad)
}

func (a *App) discoverMaster() {
	max_timeout := a.Options.MaxTimeout
	max_connect_try := a.Options.MaxConnectTry
	max_connect_try_timeout := a.Options.MaxConnectTryTimeout

	mt_duration := max_timeout
	mt_timer := time.NewTimer(mt_duration)
	// leader election process
	masterCandidateId := ""
	// try count of master finding
	ctry := 0

	for {
		// ignore if current process is master selected
		if a.cs.IsMaster {
			time.Sleep(time.Second)
			continue
		}

		// array to store server's id
		var ids []string
		// itreate server's address for requesting
		for _, ad := range a.ReplicasAddrs {
			// ignore current process address
			// append current process id to list
			if ad == a.Addr {
				ids = append(ids, a.Id)
				continue
			}
			// connect and request for check master to address
			fs, err := a.Client.ConnectIfMaster(ad)
			if err != nil {
				fmt.Println("failed to test : ", ad, " err=", err.Error())
				continue
			}
			// if address is available then append id to list
			if fs.Id != "" {
				ids = append(ids, fs.Id)
			}

			// if address is master or master candidate then follow master
			if fs.IsMaster || fs.Id == masterCandidateId {
				ctry = 0
				a.status = "followed"
				err := a.cs.FollowMaster(fs)
				if err != nil {
					a.status = "unknown"
					// log.Println("error an follow master : ", err.Error())
					break
				}

			}
		}
		if ctry == 0 {
			mt_timer.Reset(mt_duration)
		}

		var isTimeout, isTryExceeded bool
		select {
		case <-mt_timer.C:
			isTimeout = true

		default:
			isTryExceeded = (len(ids) > 0 && ctry >= max_connect_try)
		}

		// Leader election core algorithm
		// if master is not found then select a first server from sorted list as master
		if len(ids) == len(a.ReplicasAddrs) || isTryExceeded || isTimeout {
			ctry = 0
			// sort ids
			sort.Strings(ids)
			// select first
			masterCandidateId = ids[0]
			if masterCandidateId == a.Id {
				a.status = "master"
				a.cs.IsMaster = true
				fmt.Println("server is master")
			}
		}

		time.Sleep(max_connect_try_timeout)
		ctry++
	}

}

func (a *App) serve() {
	grpcServer := grpc.NewServer()
	a.gs = grpcServer

	listener, err := net.Listen("tcp", a.Addr)
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	proto.RegisterBroadcastServer(grpcServer, a.cs)
	grpcServer.Serve(listener)
}

func (a *App) Status() string {
	return a.status
}

func (a *App) Start() {
	a.setDefaultOptions()

	if a.Id == "" {
		id := sha256.Sum256([]byte(time.Now().String()))
		serverId := hex.EncodeToString(id[:])
		a.Id = serverId
	}

	a.status = "unknown"
	a.addServer(a.Addr)

	a.cs = a.Server.NewChatServer(a.Id)

	go a.discoverMaster()
	a.serve()

}

func (a *App) Exit() {
	if a.gs != nil {
		a.gs.Stop()
	}
}
