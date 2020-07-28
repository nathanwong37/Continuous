package api

import (
	"net"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/soheilhy/cmux"
	"github.com/temp/grpc"
	"github.com/temp/messenger"
)

func TestRun(t *testing.T) {
	conf := memberlist.DefaultLocalConfig()
	conf.Name = "Yolo"
	conf.BindAddr = "127.0.0.1"
	conf.BindPort = 8080
	conf.AdvertiseAddr = "127.0.0.1"
	conf.AdvertisePort = 8080
	test := messenger.NewMessenger(conf)
	nodes := []string{
		"localhost:7946",
	}
	test.Join(nodes)
	listen := NewListener(test)
	listen.run()

}

func TestCmux(t *testing.T) {
	conf := memberlist.DefaultLocalConfig()
	// conf.BindAddr = "127.0.0.1"
	// conf.BindPort = 4040
	// conf.AdvertiseAddr = "127.0.0.1"
	// conf.AdvertisePort = 4040
	test := messenger.NewMessenger(conf)
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
	server := grpc.NewGrpcServer(test)
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
