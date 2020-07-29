package temp

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
)

func TestHashFunc(t *testing.T) {
	conf3 := memberlist.DefaultLocalConfig()
	nodes := []string{
		"127.0.0.1:7946",
	}
	conf := memberlist.DefaultLocalConfig()
	conf.Name = "Yolo"
	conf.BindAddr = "127.0.0.1"
	conf.BindPort = 4209
	conf.AdvertiseAddr = "127.0.0.1"
	conf.AdvertisePort = 4209
	conf2 := memberlist.DefaultLocalConfig()
	conf2.Name = "NotFeelingLucky"
	conf2.BindAddr = "127.0.0.1"
	conf2.BindPort = 2134
	conf2.AdvertiseAddr = "127.0.0.1"
	conf2.AdvertisePort = 2134
	conf4 := memberlist.DefaultLocalConfig()
	conf4.Name = "Test1"
	conf4.BindAddr = "127.0.0.1"
	conf4.BindPort = 6969
	conf4.AdvertiseAddr = "127.0.0.1"
	conf4.AdvertisePort = 6969
	test := NewMessenger(conf3)
	test.Join(nodes)
	test2 := NewMessenger(conf)
	test2.Join(nodes)
	test3 := NewMessenger(conf2)
	test3.Join(nodes)
	test4 := NewMessenger(conf4)
	test4.Join(nodes)
	time.Sleep(5 * time.Second)
	test.printKeys()
	// a := test.syncShards()
	// b := test2.syncShards()
	// c := test3.syncShards()
	// d := test4.syncShards()
	test2.shutDown()
	test.nodeLeave(test2.M.LocalNode().Address())
	time.Sleep(10 * time.Second)
	test.printKeys()
	// a1 := cogs.NewDirector()
	// b1 := cogs.NewDirector()
	// c1 := cogs.NewDirector()
	// d1 := cogs.NewDirector()

	// a1.UpdateShards(a)
	// b1.UpdateShards(b)
	// c1.UpdateShards(c)
	// d1.UpdateShards(d)
}

func TestDirector(t *testing.T) {
	conf3 := memberlist.DefaultLocalConfig()
	nodes := []string{
		"127.0.0.1:7946",
	}
	test := NewMessenger(conf3)
	test.Join(nodes)
	a := test.syncShards()
	b := NewDirector()
	b.UpdateShards(a)

}

func TestBinarySearch(t *testing.T) {
	arrs := [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	var arr []int = arrs[0:9]
	a := binarySearch(arr, 0, len(arr)-1, 7)
	fmt.Printf("%d\n", a)
}

func TestAddress(t *testing.T) {
	uuid := "ac65b87e-d5b6-4131-b4e2-789d1fc98b7a"
	conf := memberlist.DefaultLocalConfig()
	nodes := []string{
		"127.0.0.1:7946",
	}
	test := NewMessenger(conf)
	test.Join(nodes)
	conf2 := memberlist.DefaultLocalConfig()
	conf2.Name = "NotFeelingLucky"
	conf2.BindAddr = "127.0.0.1"
	conf2.BindPort = 2134
	conf2.AdvertiseAddr = "127.0.0.1"
	conf2.AdvertisePort = 2134
	test2 := NewMessenger(conf2)
	test2.Join(nodes)
	time.Sleep(1 * time.Second)
	a, b := test.GetAddress(uuid)
	fmt.Println(a)
	fmt.Printf("%d\n", b)
}
func (messenger *Messenger) printKeys() {
	for k := range messenger.keys {
		fmt.Printf("%d  %d\n", messenger.keys[k], k)
	}
}
func (messenger *Messenger) waitForChannel() {
	messenger.ReadFromChannel()
}
