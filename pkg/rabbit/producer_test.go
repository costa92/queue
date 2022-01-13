package rabbit

import (
	"testing"
)

func TestProducer(t *testing.T) {
	options := NewRabbitmqOptions()
	rabbit := NewRabbitmq(options)
	queueInfo1 := NewQueue("23423424", "exchange_queues_ml_headline_gift_msg", "route_queues_ml_headline_gift_msg", "topic")
	prod := &ProduceTest{queueInfo1}
	rabbit.RegisterProducer(prod)

	// 第二个
	queueInfo := NewQueue("queue_heritage", "exchange_heritage", "route_heritage", "topic")
	prods := &ProduceTest2{queueInfo}

	//rabbit.listenProducer(prods)
	rabbit.RegisterProducer(prods)
	rabbit.ProducerStart()

}

type ProduceTest struct {
	*QueueOption
}

func (p *ProduceTest) Push() string {
	return "4234234"
}



type ProduceTest2 struct {
	*QueueOption
}

func (p *ProduceTest2) Push() string {
	return "3123123131312312"
}
