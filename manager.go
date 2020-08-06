package temp

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	proto "github.com/temp/plugins"
)

//Manager struct is used in order to deal with Workers/Timers within the Shard
//Remember add locks and unlocks
type Manager struct {
	shardID int
	worker  map[string]*Worker
}

//NewManager is a constructor for Manager
func NewManager(ID int) *Manager {
	return &Manager{shardID: ID, worker: make(map[string]*Worker)}
}

//CreateTimer should create a worker to map the timers uuid to
func (manager *Manager) CreateTimer(Info *proto.TimerInfo) error {
	uuid := Info.GetTimerID()
	newWork := NewWorker(Info, manager)
	manager.worker[uuid] = newWork
	start := Info.GetStartTime()
	if Info.GetMostRecent() != "" {
		start = Info.GetMostRecent()
	}
	err := newWork.CallBack(start, newWork.RunTimer)
	if err != nil {
		delete(manager.worker, uuid)
	}
	return err
}

//DeleteEntry is used to delete the entry in the worker as well as remove entry from database
func (manager *Manager) DeleteEntry(uuidstr, namespace string) {
	//removes from database
	manager.remove(uuidstr, namespace)
	//deletes
	delete(manager.worker, uuidstr)
}

//remove tries to delete the entry from the databse
func (manager *Manager) remove(uuidstr, namespace string) {
	uu, err := uuid.Parse(uuidstr)
	if err != nil {
		return
	}
	transporter := Transport{}
	res, errs := transporter.Remove(uu, namespace)
	if errs != nil {
		return
	}
	if res {
		fmt.Println(uuidstr + " Removed")
	}
}

//PullAllTimers is for manager to pull all the timers it owns from database if it has not done so
func (manager *Manager) PullAllTimers() {
	transport := new(Transport)
	timers, err := transport.GetRows(manager.shardID)
	if err != nil {
		fmt.Println("Failed to get rows")
		return
	}
	for i := range timers {
		//if there already exists a value, assume that it is the correct one
		if _, ok := manager.worker[timers[i].GetTimerID()]; ok {
			continue
		}
		manager.CreateTimer(timers[i])
	}
}

//DeleteTime will send the stop signal to the go routine
func (manager *Manager) DeleteTime(uuidstr string) bool {
	manager.worker[uuidstr].SendSignal(true)
	var count int = 0
	for count <= 3 {
		if _, ok := manager.worker[uuidstr]; ok {
			count++
			time.Sleep(10 * time.Millisecond)
			continue
		}
		break
	}
	if count > 3 {
		return false
	}
	return true
}

//StopAllTimers stops all the timers
func (manager *Manager) StopAllTimers() {
	if len(manager.worker) == 0 {
		return
	}
	for _, worker := range manager.worker {
		worker.SendSignal(false)
	}
}
