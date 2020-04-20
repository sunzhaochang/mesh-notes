package main

import (
	rls "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v2"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// server is used to implement rls.RateLimitService
type server struct {
	// limit specifies if the next request is to be rate limited
	limit bool
}

func (s *server) ShouldRateLimit(ctx context.Context,
	request *rls.RateLimitRequest) (*rls.RateLimitResponse, error) {
	log.Printf("request: %v\n", request)

	// logic to rate limit every second request
	var overallCode rls.RateLimitResponse_Code
	if s.limit {
		overallCode = rls.RateLimitResponse_OVER_LIMIT
		s.limit = false
	} else {
		overallCode = rls.RateLimitResponse_OK
		s.limit = true
	}

	response := &rls.RateLimitResponse{OverallCode: overallCode}
	log.Printf("response: %v\n", response)
	return response, nil
}

func main() {
	// create a TCP listener on port 8081
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	// create a gRPC server and register the RateLimitService server
	s := grpc.NewServer()
	rls.RegisterRateLimitServiceServer(s, &server{limit: false})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
