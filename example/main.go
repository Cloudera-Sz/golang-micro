package main

import (
	"context"
	"fmt"
	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	"os"
	"time"
)

func sayHello(s string) {
	fmt.Println("hello " + s)
}

func main() {
	if err := os.Setenv("ETCD_SERVER", "192.168.1.52:2379"); err != nil {
		panic(err.Error())
	}
	client, err := etcd.NewClient(time.Duration(5) * time.Second)
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	_, err = client.Put(ctx, "test", "test_etcd")
	cancel()
	if err != nil {
		panic(err.Error())
	}
	v, err := client.Get(context.Background(), "test")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(v)
}
