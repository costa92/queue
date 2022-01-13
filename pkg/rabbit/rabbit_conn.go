package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
)

type Rabbitmq struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	receiverList []Receiver
	producerList []Producer
	url          string
	mu           sync.RWMutex
	wg           sync.WaitGroup
}

// conn 连接 rabbitmq
func (r *Rabbitmq) conn(RabbitUrl string) {
	var err error
	var mqConn *amqp.Connection
	var mqChan *amqp.Channel
	mqConn, err = amqp.Dial(RabbitUrl)
	r.connection = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开链接失败:%s \n", err)
	}
	mqChan, err = mqConn.Channel()
	r.channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开管道失败:%s \n", err)
	}
}

func NewRabbitmq(options *RabbitmqOptions) *Rabbitmq {
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		options.Username, options.Password, options.Host, options.Port, options.Vhost)
	rabbit := &Rabbitmq{}
	rabbit.conn(rabbitUrl)
	rabbit.url = rabbitUrl
	return rabbit
}

// Close 关闭rabbit
func (r *Rabbitmq) Close() {
	if r.channel != nil || r.connection != nil {
		if err := r.channel.Close(); err != nil {
			fmt.Printf("MQ管道关闭失败：%s \n", err)
		}
		if err := r.connection.Close(); err != nil {
			fmt.Printf("MQ链接关闭失败:%s \n", err)
		}
	}
}
