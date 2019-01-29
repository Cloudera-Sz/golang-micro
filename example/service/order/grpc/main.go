package main

import (
	"context"
	"fmt"
	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	"github.com/Cloudera-Sz/golang-micro/clients/rabbitmq"
	"github.com/Cloudera-Sz/golang-micro/example/proto/base"
	"github.com/Cloudera-Sz/golang-micro/example/service/order/proto"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
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

	topic := "microtest"
	mqCli, err := rabbitmq.NewClientFromEtcd(etcdCli, "order", "dev", []string{topic})
	body := string(in.Order.Id)
	err = mqCli.PublishMessage(ctx, "", topic, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	//body := string(in.Order.Id)
	//err = mqCli.Chs[topic].Publish(
	//	"",    // exchange
	//	topic, // routing key
	//	false, // mandatory
	//	false, // immediate
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(body),
	//	})
	log.Printf(" [x] Sent %s", body)
	rabbitmq.FailOnError(err, "Failed to publish a message")

	return &pb.OrderCreateResponse{Error: &base.Error{Code: 0, Message: "OK"}}, nil
}

var etcdCli *etcd.Client
var err error

func main() {
	etcdCli, err = etcd.NewClient(5*time.Second, "192.168.1.52:2379")
	if err != nil {
		log.Panicln(err)
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err.Error())
	}
	s := grpc.NewServer()

	//mg.RegisterService(etcdCli, "order", "dev", ":50051", 20)
	//mg.SignalHandler(s, func() {
	//	fmt.Println("------------------")
	//})
	pb.RegisterOrderServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		//log.Fatalf("failed to serve: %v", err)
		fmt.Println(err.Error())
	}

}
