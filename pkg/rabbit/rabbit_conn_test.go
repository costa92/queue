package rabbit

import (
	"errors"
	"fmt"
	"testing"
)

func TestConn(t *testing.T) {

	options := NewRabbitmqOptions()
	rabbit := NewRabbitmq(options)

	queueInfo1 := NewQueue("23423424", "exchange_queues_ml_headline_gift_msg", "route_queues_ml_headline_gift_msg", "topic")
	prod := &Prod{queueInfo1}
	rabbit.RegisterReceiver(prod)

	queueInfo := NewQueue("queue_heritage", "exchange_heritage", "route_heritage", "topic")
	prods := &Prods{queueInfo}
	// 添加监听
	rabbit.RegisterReceiver(prods)

	rabbit.ReceiverStart()

}

type Prod struct {
	*QueueOption
}

func (p *Prod) Consumer(data []byte) error {
	fmt.Println(string(data))
	return errors.New("111")
}

func (p *Prod) ErrorNotify(data []byte) error {
	fmt.Println(string(data))
	fmt.Println("错误处理")
	return nil
}

type Prods struct {
	*QueueOption
}

func (p *Prods) Consumer(data []byte) error {
	fmt.Println(string(data))
	return nil
}
func (p *Prods) ErrorNotify(data []byte) error {
	fmt.Println(string(data))
	return nil
}
