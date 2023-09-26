package main

import (
	"api-channel/pkg/app"
	"api-channel/pkg/grpc_server"
	"api-channel/pkg/helper"
	"flag"
	"time"
)

func main() {

	port := flag.String("p", "8082", "The port of server")
	replica := flag.String("replicas", "", "The replication address")
	flag.Parse()

	hostPort := "localhost:" + *port
	addrs := helper.SplitArgs(*replica)

	app := &app.App{
		Addr:          hostPort,
		ReplicasAddrs: addrs,
		Server:        &grpc_server.GRPCServer{},
		Client:        &grpc_server.GRPCClient{},

		Options: &app.AppOptions{
			MaxTimeout: time.Second * 10,
		},
	}

	app.Start()

}
