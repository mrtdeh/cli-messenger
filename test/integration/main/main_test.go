package main_test

import (
	"api-channel/pkg/app"
	"api-channel/pkg/grpc_server"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestStart(t *testing.T) {

	var servers = []struct {
		Id       string
		Addr     string
		RplAddrs []string
	}{
		{
			Id:       "test-1",
			Addr:     "localhost:8081",
			RplAddrs: []string{"localhost:8081", "localhost:8082", "localhost:8083"},
		},
		{
			Id:       "test-2",
			Addr:     "localhost:8082",
			RplAddrs: []string{"localhost:8081", "localhost:8082", "localhost:8083"},
		},
		{
			Id:       "test-3",
			Addr:     "localhost:8083",
			RplAddrs: []string{"localhost:8081", "localhost:8082"},
		},
	}
	var wg sync.WaitGroup
	for _, s := range servers {
		wg.Add(1)
		a := app.App{
			Id:            s.Id,
			Addr:          s.Addr,
			ReplicasAddrs: s.RplAddrs,
			Server:        &grpc_server.GRPCServer{},
			Client:        &grpc_server.GRPCClient{},
		}
		go func() {
			defer wg.Done()
			a.Start()
		}()
		go func() {
			time.Sleep(time.Second * 5)
			for {
				if st := a.Status(); st != "" && st != "unknown" {
					fmt.Println(a.Id, a.Status())
					a.Exit()
					break
				}

			}
		}()

	}

	wg.Wait()
}
