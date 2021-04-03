package hubs

import (
	"webhook/pkg/contracts"
	"webhook/pkg/infrastructure"
)

type SendAnakinHub struct {
	rabbitMqClient *infrastructure.RabbitMqClient
}

func NewSendAnakinHub(rabbitMqClient *infrastructure.RabbitMqClient) *SendAnakinHub {
	return &SendAnakinHub{
		rabbitMqClient: rabbitMqClient,
	}
}

func (hub *SendAnakinHub) Transmit(message *contracts.BridgeMessageContract, hubConfiguration *contracts.HubConfiguration) error {
	return nil
}
