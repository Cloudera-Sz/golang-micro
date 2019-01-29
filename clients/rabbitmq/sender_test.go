package rabbitmq

import (
	"log"
	"testing"
	"time"

	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	"github.com/streadway/amqp"
)

func TestSendHello(t *testing.T) {
	etcdCli, err := etcd.NewClient(5*time.Second, "")
	if err != nil {
		log.Panicln(err)
	}
	topic := "hello2"
	mqCli, err := NewClientFromEtcd(etcdCli, "user", "dev", []string{topic})
	body := "hello"
	err = mqCli.Chs[topic].Publish(
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	FailOnError(err, "Failed to publish a message")
}
