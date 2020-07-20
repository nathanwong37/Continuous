package grpc

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestGrpcConnectionServer(t *testing.T) {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	server := new(GrpcServer)
	server.Serve(listener)
	//don't want to close server right away
	time.Sleep(10 * time.Second)
}

func TestClientConnection(t *testing.T) {
	fmt.Println("Testing")
	client := NewGrpcClient(nil)
	// _, err := client.Connect("localhost:4040")
	// if err != nil {
	// 	panic(err)
	// }

	_, err := client.CreateTimer(5, "Nathan Wong", "00:00:10", "2020")
	if err != nil {
		panic(err)
	}

}