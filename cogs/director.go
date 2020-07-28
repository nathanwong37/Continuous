package cogs

import (
	"sort"
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
		if shardInt[index2] < director.managers[index].shardID {
			manager := NewManager(shardInt[index2])
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
		manager := NewManager(shardInt[index2])
		updateManager = append(updateManager, manager)
		index2++
	}
	director.managers = updateManager

}
