package hubs

import (
	"webhook/pkg/contracts"
	"webhook/pkg/infrastructure"
)

type SendNetmonHub struct {
	rabbitMqClient *infrastructure.RabbitMqClient
}

func NewSendNetmonHub(rabbitMqClient *infrastructure.RabbitMqClient) *SendNetmonHub {
	return &SendNetmonHub{
		rabbitMqClient: rabbitMqClient,
	}
}

func (hub *SendNetmonHub) Transmit(message *contracts.BridgeMessageContract, hubConfiguration *contracts.HubConfiguration) error {
	return nil
}
