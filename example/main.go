package main

import (
	"context"
	"fmt"
	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	"github.com/Cloudera-Sz/golang-micro/clients/gorm"
	"github.com/Cloudera-Sz/golang-micro/example/service/coupon"
	"github.com/Cloudera-Sz/golang-micro/example/service/order"
	"github.com/Cloudera-Sz/golang-micro/example/service/user"
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
	client, err := etcd.NewClient(time.Duration(5)*time.Second, "")
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	_, err = client.Put(ctx, "/test", "test_etcd")
	cancel()
	if err != nil {
		panic(err.Error())
	}
	v, err := client.Get(context.Background(), "/ttt")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(v)

	db, err := gorm.NewClientFromEtcd(client, "user", "dev", &order.Order{}, &user.User{}, &coupon.Coupon{})
	if err != nil {
		panic(err.Error())
	}
	db.ClientWithContext(context.Background(), "hello")
	db.Model(&user.User{Id: 1, UserName: "admin"}).Update(&user.User{Id: 1})
	time.Sleep(time.Duration(20) * time.Second)
}
