package rabbitmq

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
)

func TestReceiveHello(t *testing.T) {
	etcdCli, err := etcd.NewClient(5*time.Second, "")
	if err != nil {
		log.Panicln(err)
	}
	topic := "hello2"
	mqCli, err := NewClientFromEtcd(etcdCli, "user", "dev", []string{topic})

	msgs, err := mqCli.Chs[topic].Consume(
		topic, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Println(d.Headers)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
