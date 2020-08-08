package temp

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
)

func TestRun(t *testing.T) {
	// conf := memberlist.DefaultLocalConfig()
	// test := NewMessenger(conf)
	nodes := []string{
		"localhost:8310",
	}
	// test.Join(nodes)
	conf2 := memberlist.DefaultWANConfig()
	// conf2.Name = "NotFeelingLucky"
	conf2.BindPort = 2134
	conf2.AdvertisePort = 2134
	test2 := NewMessenger(conf2)
	fmt.Println(test2.M.LocalNode().Address())
	test2.Join(nodes)
	// time.Sleep(2 * time.Second)
	// test2.shutDown()
	// fmt.Println("TEST2 SHUTDOWN")
	time.Sleep(25 * time.Second)
}

func (messenge *Messenger) printShard() {
	for i := range messenge.director.managers {
		if messenge.director.managers[i].shardID == 489 || messenge.director.managers[i].shardID == 59 {
			fmt.Println(messenge.M.LocalNode())
			fmt.Printf("%d\n %d \n", messenge.director.managers[i].shardID, i)
		}
	}
}
