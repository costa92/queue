package rabbit

import (
	"fmt"
)

// Receiver 定义接收者接口
type Receiver interface {
	Consumer([]byte) error    // 消费逻辑
	ErrorNotify([]byte) error // 错误回调逻辑
	QueueName() string
	ExcName() string
	QueueTopic() string
	RtName() string
}

// RegisterReceiver 注册消费信息
func (r *Rabbitmq) RegisterReceiver(receiver Receiver) {
	r.mu.Lock()
	r.receiverList = append(r.receiverList, receiver)
	r.mu.Unlock()
}

// listenReceiver 监听消费者
func (r *Rabbitmq) listenReceiver(receiver Receiver) {
	defer r.Close()
	_, err := r.channel.QueueDeclarePassive(receiver.QueueName(), true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.channel.QueueDeclare(receiver.QueueName(), true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}
	}
	// 绑定任务
	err = r.channel.QueueBind(receiver.QueueName(), receiver.RtName(), receiver.ExcName(), true, nil)
	if err != nil {
		fmt.Printf("绑定队列失败:%s \n", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.channel.Qos(1, 0, true)
	msgList, err := r.channel.Consume(receiver.QueueName(), "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("获取消费通道异常:%s \n", err)
		return
	}
	for msg := range msgList {
		if err := receiver.Consumer(msg.Body); err != nil { //处理信息
			if err := receiver.ErrorNotify(msg.Body); err != nil { // 处理错误通知
				fmt.Printf("处理错误信息错误:%s \n", err)
			}
			if err = msg.Ack(true); err != nil {
				fmt.Printf("确认消息未完成异常:%s \n", err)
				return
			}
		} else {
			if err := msg.Ack(false); err != nil { // 手动 ack
				fmt.Printf("确认消息未完成异常:%s \n", err)
				return
			}
		}
	}
}

// ReceiverStart 消费开始
func (r *Rabbitmq) ReceiverStart() {
	defer r.wg.Done()
	if len(r.receiverList) > 0 {
		for _, receiver := range r.receiverList {
			r.wg.Add(1)
			go r.listenReceiver(receiver)
		}
	}
	r.wg.Wait()
}
