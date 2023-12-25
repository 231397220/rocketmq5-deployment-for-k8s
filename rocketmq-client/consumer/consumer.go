// 运行方法:go run consumer.go -n 10.43.211.14:9876 -t sam-test1 -g sam-consumer-group
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	var nameServer, topic, group string

	flag.StringVar(&nameServer, "n", "127.0.0.1:9876", "NameServer 地址")
	flag.StringVar(&topic, "t", "testTopic", "消费主题")
	flag.StringVar(&group, "g", "myConsumerGroup", "消费者组")
	flag.Parse()

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(strings.Split(nameServer, ",")),
		consumer.WithGroupName(group),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		fmt.Printf("create consumer error: %s\n", err.Error())
		os.Exit(1)
	}

	err = c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
		fmt.Printf("Received message: Topic=%s, Body=%s, Offset=%d\n", msg.Topic, string(msg.Body), msg.QueueOffset)
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		fmt.Printf("subscribe error: %s\n", err.Error())
		os.Exit(1)
	}

	err = c.Start()
	if err != nil {
		fmt.Printf("start consumer error: %s\n", err.Error())
		os.Exit(1)
	}

	// 为了保持程序运行，您可能需要添加逻辑来避免程序直接退出
	select {} // 无限等待，防止主 goroutine 退出
}
