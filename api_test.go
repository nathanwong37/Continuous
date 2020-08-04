package temp

import (
	"net"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/soheilhy/cmux"
	//"github.com/temp/grpc"
	//"github.com/temp/messenger"
)

func TestRun(t *testing.T) {
	conf := memberlist.DefaultLocalConfig()
	test := NewMessenger(conf)
	nodes := []string{
		"localhost:7946",
	}
	test.Join(nodes)
	// conf2 := memberlist.DefaultLocalConfig()
	// conf2.Name = "NotFeelingLucky"
	// conf2.BindAddr = "127.0.1.1"
	// conf2.BindPort = 2134
	// conf2.AdvertiseAddr = "127.0.0.1"
	// conf2.AdvertisePort = 2134
	// test2 := NewMessenger(conf2)
	// test2.Join(nodes)
	//	time.Sleep(2 * time.Second)
	time.Sleep(25 * time.Second)

}

func TestCmux(t *testing.T) {
	conf := memberlist.DefaultLocalConfig()
	// conf.BindAddr = "127.0.0.1"
	// conf.BindPort = 4040
	// conf.AdvertiseAddr = "127.0.0.1"
	// conf.AdvertisePort = 4040
	test := NewMessenger(conf)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	nodes := []string{
		"localhost:7946",
	}
	test.Join(nodes)
	tcpm := cmux.New(l)
	httpl := tcpm.Match(cmux.HTTP1())
	grpcl := tcpm.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	server := NewGrpcServer(test)
	listen := NewListener(test)
	server.Serve(grpcl)
	listen.runlisten(httpl)
	go func() {
		if err := tcpm.Serve(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(5 * time.Second)

}
