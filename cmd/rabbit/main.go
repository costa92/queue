package main

import (
	"errors"
	"fmt"
	"github.com/costa92/queue/pkg/rabbit"
)

func main() {
	// config
	options := rabbit.NewRabbitmqOptions()
	rabbitServer := rabbit.NewRabbitmq(options)
	//消费者
	receiver(rabbitServer)
}

func receiver(rabbitServer *rabbit.Rabbitmq) {
	// 第一个
	queueInfo1 := rabbit.NewQueue("queuename", "exchange_queues_ml_headline_gift_msg", "route_queues_ml_headline_gift_msg", "topic")
	prod := &Prod{queueInfo1}
	rabbitServer.RegisterReceiver(prod) // 注册服务消息

	// 第二个
	queueInfo := rabbit.NewQueue("queue_heritage", "exchange_heritage", "route_heritage", "topic")
	prods := &Prods{queueInfo}
	rabbitServer.RegisterReceiver(prods)

	// 消费者执行
	rabbitServer.ReceiverStart()
}

type Prod struct {
	*rabbit.QueueOption
}

// Consumer 消费者逻辑
func (p *Prod) Consumer(data []byte) error {
	fmt.Println(string(data))
	return errors.New("111")
}

// ErrorNotify 错误信息通知
func (p *Prod) ErrorNotify(data []byte) error {
	fmt.Println(string(data))
	fmt.Println("错误处理")
	return nil
}

// Prods c
type Prods struct {
	*rabbit.QueueOption
}

func (p *Prods) Consumer(data []byte) error {
	fmt.Println(string(data))
	return nil
}
func (p *Prods) ErrorNotify(data []byte) error {
	fmt.Println(string(data))
	return nil
}
