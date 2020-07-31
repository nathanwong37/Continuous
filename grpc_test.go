package temp

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
	//"github.com/temp/messenger"
	//"github.com/temp/messenger"
)

func TestGrpcConnectionServer(t *testing.T) {
	conf := memberlist.DefaultLocalConfig()
	test := NewMessenger(conf)
	nodes := []string{
		"localhost:7946",
	}
	test.Join(nodes)
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	server := NewGrpcServer(test)
	server.Serve(listener)
	//don't want to close server right away
	time.Sleep(10 * time.Second)
}

func TestClientConnection(t *testing.T) {
	fmt.Println("Testing")
	conf := memberlist.DefaultLocalConfig()
	test := NewMessenger(conf)
	client := NewGrpcClient(nil, test)
	_, err := client.Connect("localhost:4040")
	if err != nil {
		panic(err)
	}

}
