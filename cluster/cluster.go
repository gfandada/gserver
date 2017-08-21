package cluster

import (
	"github.com/gfandada/gserver/cluster/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	CLOSEF = iota // 前端流关闭
	CLOSEB        // 后端流关闭
)

// 获取一个指定服务的路由流
func GetRouterStream(service string) pb.ClusterService_RouterClient {
	conn := GetService(service)
	if conn == nil {
		return nil
	}
	clusterServiceClient := pb.NewClusterServiceClient(conn)
	context := metadata.NewContext(context.Background(),
		metadata.New(map[string]string{"userid": ""}))
	stream, errs := clusterServiceClient.Router(context)
	if errs != nil {
		return nil
	}
	return stream
}

// 获取所有路由服务流
func GetRouterStreams() (streams map[string]pb.ClusterService_RouterClient) {
	names := GetServiceNames()
	streams = make(map[string]pb.ClusterService_RouterClient)
	for name := range names {
		streams[name] = GetRouterStream(name)
	}
	return
}
