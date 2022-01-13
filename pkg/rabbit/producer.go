package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Producer interface {
	Push() string
	QueueName() string
	ExcName() string
	QueueTopic() string
	RtName() string
}

// 发送任务
func (r *Rabbitmq) listenProducer(producer Producer) {
	defer r.Close()
	if r.channel != nil {
		r.conn(r.url)
	}
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err := r.channel.QueueDeclarePassive(producer.QueueName(), true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.channel.QueueDeclare(producer.QueueName(), true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}
	}
	// 队列绑定
	err = r.channel.QueueBind(producer.QueueName(), producer.RtName(), producer.ExcName(), true, nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
		return
	}


	// 用于检查交换机是否存在,已经存在不需要重复声明
	err = r.channel.ExchangeDeclarePassive(producer.ExcName(), producer.QueueTopic(), true, false, false, true, nil)

	if err != nil {
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err = r.channel.ExchangeDeclare(producer.ExcName(), producer.QueueTopic(), true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册交换机失败:%s \n", err)
			return
		}
	}
	fmt.Println(producer.ExcName())
	// 发送任务消息
	err = r.channel.Publish(producer.ExcName(), producer.RtName(), false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(producer.Push()),
	})
	fmt.Println(producer.Push())
	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
		return
	}
}

// RegisterProducer 注册发送指定队列指定路由的生产者
func (r *Rabbitmq) RegisterProducer(producer Producer) {
	r.mu.Lock()
	r.producerList = append(r.producerList, producer)
	r.mu.Unlock()
}

func (r *Rabbitmq) ProducerStart() {
	//defer r.wg.Done()
	if len(r.producerList) > 0 {
		for _, producer := range r.producerList {
			//r.wg.Add(1)
			r.listenProducer(producer)
		}
		//r.wg.Wait()
	}
}
