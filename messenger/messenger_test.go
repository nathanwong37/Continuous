package messenger

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
	test2.printKeys()
	test3.printKeys()
	test4.printKeys()
}

func (messenger *Messenger) printKeys() {
	for k := range messenger.keys {
		fmt.Printf("%d  %d\n", messenger.keys[k], k)
	}
}
func (messenger *Messenger) waitForChannel() {
	messenger.ReadFromChannel()
}
