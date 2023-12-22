// 运行方法:go run consumer.go -n 10.43.211.14:9876 -t sam-test1 -g sam-consumer-group
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

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
		consumer.WithNameServer([]string{nameServer}),
		consumer.WithGroupName(group),
		consumer.WithConsumerModel(consumer.Clustering),
	)
	if err != nil {
		fmt.Printf("create consumer error: %s\n", err.Error())
		os.Exit(1)
	}

	err = c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			// 打印消息内容到屏幕
			fmt.Printf("Received message: Topic=%s, Body=%s\n", msgs[i].Topic, string(msgs[i].Body))
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		fmt.Printf("subscribe error: %s\n", err.Error())
		os.Exit(1)
	}

	// 启动消费者
	err = c.Start()
	if err != nil {
		fmt.Printf("start consumer error: %s\n", err.Error())
		os.Exit(1)
	}
	time.Sleep(time.Minute * 10) // 保持消费者运行一段时间
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown consumer error: %s\n", err.Error())
	}
}
