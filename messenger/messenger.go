package messenger

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/temp/cogs"
	proto "github.com/temp/plugins"
)

//Messenger contains memberlist, shard amount and partition amount for the members
type Messenger struct {
	M         *memberlist.Memberlist
	channel   chan memberlist.NodeEvent
	shard     int
	partition int
	hashring  map[int]string
	lock      *sync.RWMutex
	keys      []int
	director  *cogs.Director
}

//NewMessenger is a constructor for messenger. It first creates a memberlist of its own, then should attempt to join
func NewMessenger(conf *memberlist.Config) *Messenger {
	if conf == nil {
		conf = memberlist.DefaultLocalConfig()
	}
	ch := make(chan memberlist.NodeEvent, 3)
	conf.Events = &memberlist.ChannelEventDelegate{ch}
	list, err := memberlist.Create(conf)
	if err != nil {
		panic(err)
	}
	//readFromChannel(ch)
	return &Messenger{
		M:         list,
		channel:   ch,
		shard:     1000,
		partition: 10,
		hashring:  make(map[int]string),
		lock:      new(sync.RWMutex),
		keys:      make([]int, 0, 1000),
		director:  cogs.NewDirector(),
	}
}

//Join functionality is used to try and join other memberlists
func (messenger *Messenger) Join(addr []string) int {
	go messenger.ReadFromChannel()
	try, err := messenger.M.Join(addr)
	if err != nil {
		return -1
	}
	return try
}

func hash(toHash string, shard int) int {
	hashfunc := sha256.New()
	hashfunc.Write([]byte(toHash))
	num := new(big.Int)
	num.SetBytes(hashfunc.Sum(nil))
	shardAmt := new(big.Int)
	shardAmt.SetInt64(int64(shard))
	num = num.Mod(num, shardAmt)
	return int(num.Int64())
}

func (messenger *Messenger) shutDown() {
	messenger.M.Leave(time.Duration(1) * time.Second)
	messenger.M.Shutdown()
}

//ReadFromChannel is for messenger to read
func (messenger *Messenger) ReadFromChannel() {
	for {
		event := <-messenger.channel
		switch ok := event.Event; ok {
		//3 Cases
		//0. Node Joins
		//1. Node leaves
		//2. Node updates
		case 0:
			messenger.lock.Lock()
			messenger.nodeJoin(event.Node.Address())
			messenger.lock.Unlock()
			shards := messenger.syncShards()
			//now messenger will command director to create shards and do all that
			messenger.director.UpdateShards(shards)
			//Need something to sync up the timers
		case 1:
			messenger.lock.Lock()
			messenger.nodeLeave(event.Node.Address())
			messenger.lock.Unlock()
			shards := messenger.syncShards()
			messenger.director.UpdateShards(shards)
		case 2:
			fmt.Println("Something updated")
		default:
			fmt.Println("Error and default")
		}
	}
}

//gives a map of the starting and ending points of the shards this node owns
func (messenger *Messenger) syncShards() map[int]int {
	shards := make(map[int]int)
	messenger.lock.Lock()
	for key := 1; key < len(messenger.keys); key++ {
		if strings.EqualFold(messenger.hashring[messenger.keys[key-1]], messenger.M.LocalNode().Address()) {
			if key-1 == 0 {
				shards[0] = messenger.keys[key-1]
			}
			shards[messenger.keys[key-1]] = messenger.keys[key]
		}
		if key == len(messenger.keys)-1 {
			if strings.EqualFold(messenger.hashring[messenger.keys[key]], messenger.M.LocalNode().Address()) {
				shards[messenger.keys[key]] = messenger.shard - 1
			}
		}
	}
	messenger.lock.Unlock()
	return shards
}

//nodeJoin should be called when a node Joins
func (messenger *Messenger) nodeJoin(address string) {
	for i := 0; i < messenger.partition; i++ {
		ans := hash(address+"."+string(i), messenger.shard)
		//if ans is in hashring... potential bug with this
		//Nodes become aware of each other at differnt times, if they hash to the same spot
		//Both Node will think the other owns the node at ans + 1, so those timers don't run
		//potential fix don't override have both nodes run timers
		for _, ok := messenger.hashring[ans]; ok; _, ok = messenger.hashring[ans] {
			ans++
			if ans > 999 {
				ans = 0
			}
		}
		messenger.hashring[ans] = address
		messenger.keys = append(messenger.keys, ans)
	}
	sort.Ints(messenger.keys)

}

//nodeLeave should be called when a node leaves
func (messenger *Messenger) nodeLeave(address string) {
	for key, value := range messenger.hashring {
		if value == address {
			//delete from key
			target := key
			index := binarySearch(messenger.keys, 0, len(messenger.keys)-1, target)
			messenger.keys = append(messenger.keys[:index], messenger.keys[index+1:]...)
			//delete from map
			delete(messenger.hashring, key)
		}
	}
}

//recursive binarysearch
func binarySearch(search []int, low int, high int, target int) int {
	if low > high {
		return -1
	}
	mid := (low + high) / 2
	if search[mid] == target {
		return mid
	} else if search[mid] < target {
		return binarySearch(search, mid+1, high, target)
	} else {
		return binarySearch(search, low, mid-1, target)
	}

}

//GetAddress is used to get the Node address and shard of a timerId
func (messenger *Messenger) GetAddress(timerID string) (address string, shardID int) {
	//Steps
	//1. Hash the timerID to a shard
	//2. Find which timerId it belongs to
	shardNum := hash(timerID, messenger.shard)
	i := 0
	for i < len(messenger.keys) && shardNum >= messenger.keys[i] {
		i++
	}
	if i != 0 {
		i = i - 1
	}
	addr := messenger.hashring[messenger.keys[i]]
	return addr, shardNum
}

//CreateTime will be called to tell director to fire up the timer
func (messenger *Messenger) CreateTime(timerinfo *proto.TimerInfo) {
	fmt.Println("SUCKER")
}
