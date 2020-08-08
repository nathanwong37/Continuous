package temp

import (
	"context"
	"net"
	proto "github.com/temp/plugins"
	"google.golang.org/grpc"
)

//GrpcServer is a struct to call methods
type GrpcServer struct {
	messenger *Messenger
}

//NewGrpcServer is a constructor
func NewGrpcServer(messenger *Messenger) *GrpcServer {
	return &GrpcServer{
		messenger: messenger,
	}
}

//Serve is just for the server to run in the back
func (server *GrpcServer) Serve(listener net.Listener) error {
	register := grpc.NewServer()
	proto.RegisterActionsServer(register, server)
	go register.Serve(listener)
	return nil

}

//Create will do the following
//1.) Have the messenger component receive the data
//2.) returns result of the component
func (server *GrpcServer) Create(ctx context.Context, createRequest *proto.CreateJobRequest) (*proto.CreateJobResponse, error) {
	timer := proto.TimerInfo{
		TimerID:   createRequest.GetTimerId(),
		ShardID:   createRequest.GetShardId(),
		NameSpace: createRequest.GetNameSpace(),
		Interval:  createRequest.GetInterval(),
		Count:     createRequest.GetCount(),
		StartTime: createRequest.GetStartTime(),
	}
	//at this point send timerinfo struct to messenger, having messenger deal with this
	//messenger.GetTimer(&timer)
	//return the create job response now
	server.messenger.CreateTime(&timer)
	return &proto.CreateJobResponse{Timerinfo: &timer}, nil
}

//Delete will do the following
//1.) Tell the messenger what to delete
//2.) return result of messenger
func (server *GrpcServer) Delete(ctx context.Context, deleteRequest *proto.DeleteJobRequest) (*proto.DeleteJobResponse, error) {
	deleted := server.messenger.DeleteTime(deleteRequest.GetTimerId(), deleteRequest.GetNameSpace(), int(deleteRequest.GetShardId()))
	return &proto.DeleteJobResponse{Deleted: deleted}, nil
}
