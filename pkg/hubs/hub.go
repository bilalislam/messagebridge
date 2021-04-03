package hubs

import (
	"webhook/pkg/contracts"
	"webhook/pkg/infrastructure"
)

type IHub interface {
	Transmit(message *contracts.BridgeMessageContract, hubConfiguration *contracts.HubConfiguration) error
}

type Hub struct {
	hubs           map[string]IHub
	rabbitMqClient *infrastructure.RabbitMqClient
}

func NewHub() *Hub {
	rabbitMqClient := infrastructure.NewRabbitMqClient()
	return &Hub{
		hubs:           nil,
		rabbitMqClient: rabbitMqClient,
	}
}
