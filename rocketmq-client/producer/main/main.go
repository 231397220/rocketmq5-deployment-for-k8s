// 运行方法:go run main.go -n 10.43.211.14:9876 -t sam-test1

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	// 定义命令行参数
	var nameServer string
	var topic string

	// 初始化命令行参数
	flag.StringVar(&nameServer, "n", "127.0.0.1:9876", "NameServer 地址")
	flag.StringVar(&topic, "t", "testTopic", "主题")
	flag.Parse()

	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{nameServer}),
		producer.WithRetry(2),
	)
	if err != nil {
		fmt.Printf("create producer error: %s\n", err.Error())
		os.Exit(1)
	}
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s\n", err.Error())
		os.Exit(1)
	}

	// 发送消息
	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte(fmt.Sprintf("Hello RocketMQ ,this is Sam. %d", i)),
		}

		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err.Error())
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
		time.Sleep(time.Second)
	}

	// 关闭生产者
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s\n", err.Error())
	}
}
