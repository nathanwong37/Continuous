package temp

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/stretchr/testify/require"
)

func TestGrpcConnectionServer(t *testing.T) {
	conf := memberlist.DefaultLocalConfig()
	test := NewMessenger(conf)
	nodes := []string{
		"localhost:7946",
	}
	test.Join(nodes)
	_, err := test.client.CreateTimer(70, "Nathan Wong", "00:00:10", "")
	require.NoError(t, err)
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
