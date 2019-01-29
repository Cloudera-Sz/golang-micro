package main

import (
	"context"
	"fmt"
	"github.com/Cloudera-Sz/golang-micro/example/service/order/proto"
	"google.golang.org/grpc"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", CreateOrder)
	r.Run(":29090")
}

func CreateOrder(c *gin.Context) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		//log.Fatalf("did not connect: %v", err)
		fmt.Println(err.Error())
	}
	defer conn.Close()
	client := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)
	defer cancel()
	r, err := client.Create(ctx, &pb.OrderCreateRequest{})
	if err != nil {
		//log.Fatalf("could not greet: %v", err)
		fmt.Println(err.Error())
	}
	//log.Printf("Greeting: %s", r.Message)

	c.JSON(200, r.Error)
}
