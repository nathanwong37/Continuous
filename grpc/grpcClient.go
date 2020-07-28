//Might have to have a messenger be part of the client, blegh
package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/temp/messenger"
	proto "github.com/temp/plugins"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

//Client is for dial options, as well as to use important methods
type Client struct {
	dialOption []grpc.DialOption
	messenger  *messenger.Messenger
}

//NewGrpcClient is a Constructor that returns a new instance of grpc client
func NewGrpcClient(dialOpt grpc.DialOption, mess *messenger.Messenger) *Client {
	if dialOpt == nil {
		dialOpt = grpc.WithInsecure()
	}
	return &Client{
		dialOption: []grpc.DialOption{
			dialOpt,
			grpc.WithBlock(),
		},
		messenger: mess,
	}
}

//Connect function connects the client to a server
func (client *Client) Connect(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, client.dialOption...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//CreateTimer will forward the create request to the appropriate node
func (client *Client) CreateTimer(count int32, namespace, interval, startTime string) (string, error) {
	timerID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	timerIDString := timerID.String()
	addr, shardResult := client.messenger.GetAddress(timerIDString)
	shardRes := int32(shardResult)
	//bug can't connect to
	addr = "localhost:4040"
	conn, err := client.Connect(addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := proto.NewActionsClient(conn)
	resp, err := c.Create(context.Background(), &proto.CreateJobRequest{
		TimerId:   timerIDString,
		ShardId:   shardRes,
		NameSpace: namespace,
		Interval:  interval,
		Count:     count,
		StartTime: startTime,
	})
	if err != nil {
		return "", err
	}
	fmt.Println(resp.GetTimerinfo().GetTimerID())
	return resp.GetTimerinfo().GetTimerID(), nil
}

//DeleteTimer will forward the delete request to the appropriate node, assume param are authenticated
func (client *Client) DeleteTimer(uuid, namespace string) (int, error) {
	//so first we want to use the uuid to calculate the address to send
	//for now we use local host
	//addr := something.membership.getAddr(uuid)
	addr := "localhost:4040"
	conn, err := client.Connect(addr)
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	c := proto.NewActionsClient(conn)
	resp, err := c.Delete(context.Background(), &proto.DeleteJobRequest{
		TimerId:   uuid,
		NameSpace: namespace,
	})
	if err != nil {
		return -1, err
	}
	if !resp.Deleted {
		return -1, errors.New("Failed to delete")
	}
	return 1, nil
}
