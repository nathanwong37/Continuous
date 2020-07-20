package messenger

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
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
	messenger.M.Leave(time.Duration(5) * time.Second)
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
			for i := 0; i < messenger.partition; i++ {
				ans := hash(event.Node.Address()+"."+string(i), messenger.shard)
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
				messenger.hashring[ans] = event.Node.Address()
				messenger.keys = append(messenger.keys, ans)
			}
			sort.Ints(messenger.keys)
			messenger.lock.Unlock()
		//Node Left
		case 1:
			fmt.Println("Something left")
		//Node updated
		case 2:
			fmt.Println("Something updated")
		default:
			fmt.Println("Error and default")
		}
	}
}
