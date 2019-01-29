package main

import (
	"context"
	"fmt"
	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	"github.com/Cloudera-Sz/golang-micro/example/service/order/proto"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var etcdCli *etcd.Client
var err error

func main() {
	etcdCli, err = etcd.NewClient(5*time.Second, "192.168.1.52:2379")
	if err != nil {
		log.Panicln(err)
	}
	r := gin.Default()
	r.GET("/", CreateOrder)
	r.Run(":29090")
}

func CreateOrder(c *gin.Context) {
	id := c.Query("id")
	order_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		//log.Fatalf("could not greet: %v", err)
		c.JSON(200, err.Error())
		return
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		//log.Fatalf("did not connect: %v", err)
		fmt.Println(err.Error())
		c.JSON(200, err.Error())
		return
	}
	defer conn.Close()
	client := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)
	defer cancel()
	r, err := client.Create(ctx, &pb.OrderCreateRequest{Order: &pb.Order{
		Id:     order_id,
		UserId: 1,
	}})
	if err != nil {
		//log.Fatalf("could not greet: %v", err)
		fmt.Println(err.Error())
		c.JSON(200, err.Error())
		return
	}
	//log.Printf("Greeting: %s", r.Message)

	c.JSON(200, r.Error)
}
