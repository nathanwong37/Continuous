package temp

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
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

	director  *Director
	server    *GrpcServer
	client    *Client
	listen    *Listener
	transport *Transport
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
		fmt.Println(err.Error())
	}
	return &Messenger{
		M:         list,
		channel:   ch,
		shard:     1000,
		partition: 10,
		hashring:  make(map[int]string),
		lock:      new(sync.RWMutex),
		keys:      make([]int, 0, 1000),
		director:  NewDirector(),
		transport: new(Transport),
	}
}

//Join functionality is used to try and join other memberlists
//Also starts grpc server, client and api
func (messenger *Messenger) Join(addr []string) (int, error) {
	go messenger.ReadFromChannel()
	try, err := messenger.M.Join(addr)
	if err != nil {
		return -1, err
	}
	if messenger.listen == nil {
		messenger.listen = NewListener(messenger)
		go messenger.listen.run(messenger.M.LocalNode().Addr.String() + ":8080")
	}
	l, err := net.Listen("tcp", messenger.M.LocalNode().Addr.String()+":51284")
	if err != nil {
		return -1, err
	}
	if messenger.server == nil {
		messenger.server = NewGrpcServer(messenger)
	}
	messenger.server.Serve(l)
	if messenger.client == nil {
		messenger.client = NewGrpcClient(nil, messenger)
	}
	return try, nil
}

//hash function using sha256
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

//used to test messenger shutdown
func (messenger *Messenger) shutDown() {
	messenger.director.StopAllTimers()
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
		case 0:
			messenger.lock.Lock()
			messenger.nodeJoin(event.Node.Address())
			messenger.lock.Unlock()
			shards := messenger.syncShards()
			messenger.director.UpdateShards(shards, messenger.shard)
		case 1:
			messenger.lock.Lock()
			messenger.nodeLeave(event.Node.Address())
			messenger.lock.Unlock()
			shards := messenger.syncShards()
			messenger.director.UpdateShards(shards, messenger.shard)
		default:
			fmt.Println("Default")
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
		//Better fix is to probably just increase shard amount... Reduces chance of collision
		for _, ok := messenger.hashring[ans]; ok; _, ok = messenger.hashring[ans] {
			ans++
			if ans > messenger.shard-1 {
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
func (messenger *Messenger) CreateTime(timerInfo *proto.TimerInfo) {
	_, err := messenger.transport.Create(timerInfo)
	if err != nil {
		fmt.Println(err.Error())
	}
	messenger.director.CreateTimer(timerInfo)
}

//DeleteTime will be called to tell director to delete the timer
func (messenger *Messenger) DeleteTime(uuidstr, namespace string, shardid int) bool {
	return messenger.director.DeleteTime(uuidstr, namespace, shardid)
}
