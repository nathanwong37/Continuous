package temp

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
)

func TestRun(t *testing.T) {
	// conf := memberlist.DefaultLocalConfig()
	// test := NewMessenger(conf)
	nodes := []string{
		"localhost:7946",
	}
	// test.Join(nodes)
	conf2 := memberlist.DefaultLocalConfig()
	conf2.Name = "NotFeelingLucky"
	conf2.BindPort = 2134
	conf2.AdvertisePort = 2134
	conf := CustomConfig(conf2, true)
	test2 := NewMessenger(conf)
	test2.Join(nodes)
	// time.Sleep(2 * time.Second)
	// test2.shutDown()
	// fmt.Println("TEST2 SHUTDOWN")
	time.Sleep(25 * time.Second)
}

func TestIP(t *testing.T) {
	a := GetOutboundIP()
	fmt.Println(a.String())
}

func (messenge *Messenger) printShard() {
	for i := range messenge.director.managers {
		if messenge.director.managers[i].shardID == 489 || messenge.director.managers[i].shardID == 59 {
			fmt.Println(messenge.M.LocalNode())
			fmt.Printf("%d\n %d \n", messenge.director.managers[i].shardID, i)
		}
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
