package rabbit

type RabbitmqOptions struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Vhost    string `json:"v-host"`
}

func NewRabbitmqOptions() *RabbitmqOptions {
	return &RabbitmqOptions{
		Host:     "192.168.11.101",
		Port:     5672,
		Username: "dgmq-sz-bigdata",
		Password: "bigdata123456",
		Vhost:    "/queue_dfml_bigdata",
	}
}

type QueueOption struct {
	QName  string `json:"q_name"`
	ExName string `json:"ex_name"`
	RkName string `json:"rk_name"`
	Topic  string `json:"topic"`
}

func NewQueueOption() *QueueOption {
	return &QueueOption{
		QName:  "23423424",
		ExName: "exchange_queues_ml_headline_gift_msg",
		RkName: "route_queues_ml_headline_gift_msg",
		Topic:  "topic",
	}
}

func NewQueue(QName string, ExName string, RkName string, Topic string) *QueueOption {
	return &QueueOption{
		QName:  QName,
		ExName: ExName,
		RkName: RkName,
		Topic:  Topic,
	}
}

func (p *QueueOption) QueueName() string {
	return p.QName
}

func (p *QueueOption) ExcName() string {
	return p.ExName
}

func (p *QueueOption) QueueTopic() string {
	return p.Topic
}

func (p *QueueOption) RtName() string {
	return p.RkName
}
