package temp

import (
	"fmt"
	"sort"

	proto "github.com/temp/plugins"
)

//Director is in charge of Shards, also has to periodically check database
type Director struct {
	managers []*Manager
}

//NewDirector is a constructor to initialize the map
func NewDirector() *Director {
	return &Director{managers: make([]*Manager, 0)}
}

//UpdateShards add shards that are owned, and deletes shard no longer in ownership
//maps are randomized...
func (director *Director) UpdateShards(shard map[int]int) {
	var index int = 0
	var index2 int = 0
	updateManager := make([]*Manager, 0)
	var shardInt []int
	//shardInt := make([]int, 0)
	//populate with the shards we now own
	for start, end := range shard {
		for i := start; i < end; i++ {
			shardInt = append(shardInt, i)
		}
		//will think of better way for this
		if end == 999 {
			shardInt = append(shardInt, end)
		}
	}
	sort.Ints(shardInt)
	for index2 < len(shardInt) && index < len(director.managers) {
		//if... gained shards
		//else if.. same shards
		//else... lost shards
		if shardInt[index2] < director.managers[index].shardID {
			manager := NewManager(shardInt[index2])
			//Still need to pull the timers from the database
			updateManager = append(updateManager, manager)
			index2++
		} else if shardInt[index2] == director.managers[index].shardID {
			updateManager = append(updateManager, director.managers[index])
			index++
		} else {
			index++
		}
	}
	for index2 < len(shardInt) {
		//Still need to pull the timers from the database
		manager := NewManager(shardInt[index2])
		updateManager = append(updateManager, manager)
		index2++
	}
	director.managers = updateManager
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
	//d.managers[index].CreateTimer()
}

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
