package data

import (
	"sync"
)

type UserExistInfo struct {
	ServerId string
	Active   bool
}

var (
	RoomsServers = map[string][]string{}
	Users        = map[string]UserExistInfo{}
	mrw          sync.RWMutex
)

func AddServerToRoom(room, serverId string) {
	mrw.Lock()
	defer mrw.Unlock()

	if r, ok := RoomsServers[room]; ok {
		for _, s := range r {
			if s == serverId {
				return
			}
		}
		r = append(r, serverId)
		RoomsServers[room] = r
	} else {
		RoomsServers[room] = []string{serverId}
	}

}

func AddUser(name string, active bool, sid string) {
	mrw.Lock()
	defer mrw.Unlock()

	ue := UserExistInfo{
		Active:   active,
		ServerId: sid,
	}
	Users[name] = ue

}

func CleanServer(sid string) {
	for _, u := range Users {
		if u.ServerId == sid {
			u.Active = false
		}
	}
}
