// 运行方法:go run main.go -n 10.43.211.14:9876 -t sam-test1

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	var nameServer, topic string

	flag.StringVar(&nameServer, "n", "127.0.0.1:9876", "NameServer 地址")
	flag.StringVar(&topic, "t", "testTopic", "生产主题")
	flag.Parse()

	// 创建生产者实例
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{nameServer}),
		producer.WithRetry(2),
	)
	if err != nil {
		fmt.Printf("create producer error: %s\n", err.Error())
		os.Exit(1)
	}

	// 启动生产者
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s\n", err.Error())
		os.Exit(1)
	}

	// 创建一个读取器以从标准输入读取
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("RocketMQ Producer")
	fmt.Println("---------------------")
	fmt.Println("Enter 'exit' to quit.")

	for {
		// 读取用户输入
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		// 检查是否退出
		if strings.ToLower(text) == "exit" {
			break
		}

		// 创建并发送消息
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte(text),
		}

		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err.Error())
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}

	// 关闭生产者
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s\n", err.Error())
	}
}