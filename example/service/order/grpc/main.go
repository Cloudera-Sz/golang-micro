package main

import (
	"context"
	"fmt"
	"github.com/Cloudera-Sz/golang-micro/example/service/order/proto"
	"google.golang.org/grpc"
	"net"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) List(ctx context.Context, in *pb.OrderListRequest) (response *pb.OrderListResponse, err error) {
	return &pb.OrderListResponse{}, nil
}

func (s *server) Get(ctx context.Context, in *pb.OrderGetRequest) (response *pb.OrderGetResponse, err error) {
	return &pb.OrderGetResponse{}, nil
}

func (s *server) Create(ctx context.Context, in *pb.OrderCreateRequest) (response *pb.OrderCreateResponse, err error) {
	return &pb.OrderCreateResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err.Error())
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		//log.Fatalf("failed to serve: %v", err)
		fmt.Println(err.Error())
	}
}
