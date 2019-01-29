package main

import (
	"context"
	"fmt"
	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	"github.com/Cloudera-Sz/golang-micro/clients/rabbitmq"
	"log"
	"time"
)

func main() {
	etcdCli, err := etcd.NewClient(5*time.Second, "192.168.1.52:2379")
	if err != nil {
		log.Panicln(err)
	}
	topic := "microtest"
	mqCli, err := rabbitmq.NewClientFromEtcd(etcdCli, "coupon", "dev", []string{topic})

	msgs, err := mqCli.Chs[topic].Consume(
		topic, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	rabbitmq.FailOnError(err, "Failed to register a consumer")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			mqCli.ConsumeMessage(context.Background(), &d, func() {
				fmt.Println(d.Headers)
				fmt.Println("Received a message: %s", d.Body)
			})
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
