package temp

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/stretchr/testify/assert"
)

func TestHashFunc(t *testing.T) {
	conf3 := memberlist.DefaultLocalConfig()
	nodes := []string{
		"127.0.0.1:7946",
	}
	conf := memberlist.DefaultLocalConfig()
	conf.Name = "Test"
	conf.BindPort = 3264
	conf.AdvertisePort = 3264
	conf2 := memberlist.DefaultLocalConfig()
	conf2.Name = "Test2"
	conf2.BindPort = 2134
	conf2.AdvertisePort = 2134
	conf4 := memberlist.DefaultLocalConfig()
	conf4.Name = "Test3"
	conf4.BindPort = 3125
	conf4.AdvertisePort = 3125
	test := NewMessenger(CustomConfig(conf3))
	test.Join(nodes)
	test2 := NewMessenger(CustomConfig(conf))
	test2.Join(nodes)
	test3 := NewMessenger(CustomConfig(conf2))
	test3.Join(nodes)
	test4 := NewMessenger(CustomConfig(conf4))
	test4.Join(nodes)
	time.Sleep(5 * time.Second)
	test.printKeys()
	test2.printKeys()
	test3.printKeys()
	test4.printKeys()
}

func TestDirector(t *testing.T) {
	conf3 := DefaultConfig()
	nodes := []string{
		"127.0.0.1:7946",
	}
	test := NewMessenger(conf3)
	test.Join(nodes)
	a := test.syncShards()
	b := NewDirector()
	b.UpdateShards(a, test.shard)

}

func TestBinarySearch(t *testing.T) {
	arrs := [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	var arr []int = arrs[0:9]
	a := binarySearch(arr, 0, len(arr)-1, 7)
	fmt.Printf("%d\n", a)
}

func TestAddress(t *testing.T) {
	uuid := "ac65b87e-d5b6-4131-b4e2-789d1fc98b7a"
	conf := DefaultConfig()
	nodes := []string{
		"localhost:7946",
	}
	test := NewMessenger(conf)
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
