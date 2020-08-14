package continuous

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	conf "github.com/Continuous/config"
	"github.com/hashicorp/memberlist"
)

func TestRun(t *testing.T) {
	// conf := memberlist.DefaultLocalConfig()
	// test := NewMessenger(conf)
	nodes := []string{
		"127.0.0.1:2134",
	}
	// test.Join(nodes)
	config2 := memberlist.DefaultLocalConfig()
	config2.Name = "NotFeelingLucky"
	config2.BindAddr = "127.0.0.1"
	config2.BindPort = 2134
	config2.AdvertisePort = 2134
	config := conf.CustomConfig(config2, true)
	test2 := NewMessenger(config)
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
