package main

import (
	"context"
	"fmt"
	"net"

	al "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
	"google.golang.org/grpc"
)

func main() {
	exit := make(chan int)

	ctx := context.Background()
	als := &AccessLogService{}

	go RunAccessLogServer(ctx, als, 18090)

	<-exit
}

// RunAccessLogServer starts an accesslog service.
func RunAccessLogServer(ctx context.Context, als *AccessLogService, port uint) {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("listen error, err: %v", err)
	}

	al.RegisterAccessLogServiceServer(grpcServer, als)

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			fmt.Println("serve error, err: %v", err)
		}
	}()

	<-ctx.Done()

	grpcServer.GracefulStop()
}
