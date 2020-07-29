package temp

import (
	proto "github.com/temp/plugins"
)

//Manager struct is used in order to deal with Workers/Timers within the Shard
type Manager struct {
	shardID int
	worker  map[string]*Worker
}

//NewManager is a constructor for Manager
func NewManager(ID int) *Manager {
	return &Manager{shardID: ID, worker: make(map[string]*Worker)}
}

//CreateTimer should create a worker to map the timers uuid to
func (manager *Manager) CreateTimer(Info *proto.TimerInfo) {
	uuid := Info.GetTimerID()
	newWork := NewWorker(Info)
	manager.worker[uuid] = newWork
	//newWork.CreateTimer(Info)
}
