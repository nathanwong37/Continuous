package temp

import (
	"fmt"
	"sort"
	"sync"
	"time"

	proto "github.com/temp/plugins"
)

//Director is in charge of Shards, also has to periodically check database
type Director struct {
	managers []*Manager
	lock     *sync.RWMutex
}

//NewDirector is a constructor to initialize the map
func NewDirector() *Director {
	direct := &Director{
		managers: make([]*Manager, 0),
		lock:     new(sync.RWMutex),
	}
	//start periodic scanner to pick missed timers due to race conditions
	go direct.PeriodicScan()
	return direct
}

//UpdateShards add shards that are owned, and deletes shard no longer in ownership
//maps are randomized...
func (director *Director) UpdateShards(shard map[int]int, cap int) {
	var index int = 0
	var index2 int = 0
	updateManager := make([]*Manager, 0)
	var shardInt []int
	//populate with the shards we now own
	for start, end := range shard {
		for i := start; i < end; i++ {
			shardInt = append(shardInt, i)
		}
		if end == cap-1 {
			shardInt = append(shardInt, end)
		}
	}
	sort.Ints(shardInt)
	director.lock.Lock()
	defer director.lock.Unlock()
	for index2 < len(shardInt) && index < len(director.managers) {
		//if... gained shards
		//else if.. same shards
		//else... lost shards
		if shardInt[index2] < director.managers[index].shardID {
			manager := NewManager(shardInt[index2])
			updateManager = append(updateManager, manager)
			index2++
		} else if shardInt[index2] == director.managers[index].shardID {
			updateManager = append(updateManager, director.managers[index])
			index++
			index2++
		} else {
			director.StopTimer(index)
			index++
		}
	}
	for index2 < len(shardInt) {
		manager := NewManager(shardInt[index2])
		updateManager = append(updateManager, manager)
		index2++
	}
	director.managers = updateManager
	director.PullAllTimers()
}

//CreateTimer is used to create the timer in director
//binary search for manager, if manager doesn't exist create it
func (director *Director) CreateTimer(Info *proto.TimerInfo) {
	target := int(Info.GetShardID())
	index := binSearchManager(director.managers, 0, len(director.managers)-1, target)
	//create if doesn't exist
	if index == -1 {
		manager := NewManager(target)
		index = 0
		for len(director.managers) > index && target > director.managers[index].shardID {
			index++
		}
		director.managers = append(director.managers, nil)
		if index == 0 {
			copy(director.managers[1:], director.managers[0:])
		} else {
			copy(director.managers[index:], director.managers[index-1:])
		}
		director.managers[index] = manager
	}
	err := director.managers[index].CreateTimer(Info)
	if err != nil {
		fmt.Println("Error creating timer")
	}
	return
}

//binSearchManager is to search through manager in a binary search
func binSearchManager(managers []*Manager, low, high, target int) int {
	if low > high {
		return -1
	}
	mid := (low + high) / 2
	if managers[mid].shardID == target {
		return mid
	} else if managers[mid].shardID > target {
		return binSearchManager(managers, low, mid-1, target)
	} else {
		return binSearchManager(managers, mid+1, high, target)
	}
}

//PullAllTimers has the director command all managers to get timers it may have missed from database
func (director *Director) PullAllTimers() {
	for i := 0; i < len(director.managers); i++ {
		director.managers[i].PullAllTimers()
	}
}

//DeleteTime should delete the time from the manager
func (director *Director) DeleteTime(uuidstr, namespace string, shardID int) bool {
	index := binSearchManager(director.managers, 0, len(director.managers), shardID)
	if index == -1 {
		//delete request went to wrong node
		return false
	}
	return director.managers[index].DeleteTime(uuidstr)
}

//StopTimer has the director tell the manager to stop all current timers
func (director *Director) StopTimer(index int) {
	director.managers[index].StopAllTimers()
}

//StopAllTimers stops all the timers in director, should only be called when messenger shutsdowns
func (director *Director) StopAllTimers() {
	for index := range director.managers {
		director.managers[index].StopAllTimers()
	}
}

//PeriodicScan should be used to periodically scan for timers, should only be called when director is first created
func (director *Director) PeriodicScan() {
	for {
		time.Sleep(600 * time.Second)
		director.lock.Lock()
		director.PullAllTimers()
		director.lock.Unlock()
	}
}
