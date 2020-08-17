package continuous

import (
	"fmt"
	"testing"
	"time"

	conf "github.com/Continuous/config"
	"github.com/hashicorp/memberlist"
	"github.com/stretchr/testify/assert"
)

func TestMessenger(t *testing.T) {
	nodes := []string{
		"localhost:7946",
	}
	messenger1 := NewMessenger(conf.CustomConfig(memberlist.DefaultLANConfig(), true))
	messenger1.Join(nodes)
	memberconfig := memberlist.DefaultLANConfig()
	memberconfig.Name = "Node1"
	memberconfig.BindPort = 1234
	memberconfig.AdvertisePort = 1234
	messenger2 := NewMessenger(conf.CustomConfig(memberconfig, true))
	messenger2.Join(nodes)
	memberconfig2 := memberlist.DefaultLANConfig()
	memberconfig2.Name = "Node2"
	memberconfig2.BindPort = 2134
	memberconfig2.AdvertisePort = 2134
	messenger3 := NewMessenger(conf.CustomConfig(memberconfig2, true))
	messenger3.Join(nodes)
	time.Sleep(2 * time.Second)
	numNodes := messenger3.M.NumMembers()
	assert.Equal(t, 3, numNodes)
	messenger2.shutDown()
	time.Sleep(2 * time.Second)
	numNodes = messenger3.M.NumMembers()
	assert.Equal(t, 2, numNodes)
}

func TestBinarySearch(t *testing.T) {
	arrs := [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	var arr []int = arrs[0:9]
	index := binarySearch(arr, 0, len(arr)-1, 7)
	assert.Equal(t, index, 7)
}

func TestAddress(t *testing.T) {
	uuid := "ac65b87e-d5b6-4131-b4e2-789d1fc98b7a"
	config := conf.DefaultConfig()
	nodes := []string{
		"localhost:7946",
	}
	test := NewMessenger(config)
	test.Join(nodes)
	time.Sleep(1 * time.Second)
	a, b := test.GetAddress(uuid)
	assert.Equal(t, a, "192.168.5.56:7946")
	assert.Equal(t, b, 177)
}

func (messenger *Messenger) printKeys() {
	for k := range messenger.keys {
		fmt.Printf("%d  %d\n", messenger.keys[k], k)
	}
}
