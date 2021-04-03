package contracts

type HubConfiguration struct {
	BrokerConfiguration BrokerConfiguration `json:"BrokerConfiguration"`
}

type BrokerConfiguration struct {
	Nodes        []string `json:"nodes"`
	ExchangeName string   `json:"ExchangeName"`
	ExchangeType string   `json:"ExchangeType"`
	QueueName    string   `json:"QueueName"`
	Durable      bool     `json:"Durable"`
	AutoDelete   bool     `json:"AutoDelete"`
	Internal     bool     `json:"Internal"`
	NoWait       bool     `json:"NoWait"`
	RoutingKey   string   `json:"RoutingKey"`
	Mandatory    bool     `json:"Mandatory"`
	Immediate    bool     `json:"Immediate"`
	ContentType  string   `json:"Content-Type"`
	DeliveryMode string   `json:"DeliveryMode"`
}
