package continuous

import (
	"context"
	"errors"
	"strconv"
	"strings"

	proto "github.com/Continuous/plugins"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

//GrpcClient is for dial options, as well as to use important methods
type GrpcClient struct {
	dialOption []grpc.DialOption
	messenger  *Messenger
}

//NewGrpcClient is a Constructor that returns a new instance of grpc client
func NewGrpcClient(dialOpt grpc.DialOption, messenger *Messenger) *GrpcClient {
	if dialOpt == nil {
		dialOpt = grpc.WithInsecure()
	}
	return &GrpcClient{
		dialOption: []grpc.DialOption{
			dialOpt,
			grpc.WithBlock(),
		},
		messenger: messenger,
	}
}

//Connect function connects the client to a server
func (client *GrpcClient) Connect(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, client.dialOption...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//CreateTimer will forward the create request to the appropriate node
func (client *GrpcClient) CreateTimer(count int32, namespace, interval, startTime string) (string, error) {
	timerID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	timerIDString := timerID.String()
	addr, shardResult := client.messenger.GetAddress(timerIDString)
	if !client.messenger.config.LocalConnect {
		addr = trimAddress(addr)
		addr = addr + ":" + strconv.Itoa(client.messenger.config.RPCPort)
	}
	// shardRes := int32(shardResult)
	conn, err := client.Connect(addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := proto.NewActionsClient(conn)
	resp, err := c.Create(context.Background(), &proto.CreateJobRequest{
		TimerId:   timerIDString,
		ShardId:   int32(shardResult),
		NameSpace: namespace,
		Interval:  interval,
		Count:     count,
		StartTime: startTime,
	})
	if err != nil {
		return "", err
	}
	//fmt.Println(resp.GetTimerinfo().GetTimerID())
	return resp.GetTimerinfo().GetTimerID(), nil
}

//DeleteTimer will forward the delete request to the appropriate node, assume param are authenticated
func (client *GrpcClient) DeleteTimer(uuidstr, namespace string) (int, error) {
	uu, err := uuid.Parse(uuidstr)
	if err != nil {
		return -1, err
	}
	//Make sure it is actually in the database
	_, err = client.messenger.transporter.Get(uu, namespace)
	if err != nil {
		return -1, err
	}
	addr, shardResult := client.messenger.GetAddress(uuidstr)
	addr = trimAddress(addr)
	addr = addr + ":" + strconv.Itoa(client.messenger.config.RPCPort)
	conn, err := client.Connect(addr)
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	c := proto.NewActionsClient(conn)
	resp, err := c.Delete(context.Background(), &proto.DeleteJobRequest{
		TimerId:   uuidstr,
		NameSpace: namespace,
		ShardId:   int32(shardResult),
	})
	if err != nil {
		return -1, err
	}
	if !resp.Deleted {
		return -1, errors.New("Failed to delete")
	}
	return 1, nil
}

func trimAddress(s string) string {
	if index := strings.Index(s, ":"); index != -1 {
		return s[:index]
	}
	return s
}
