package hubs

import (
	"errors"
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
	rabbitMqClient := infrastructure.NewRabbitMqClient(nil)
	return &Hub{
		hubs: map[string]IHub{
			"SendNetmonHub": NewSendNetmonHub(rabbitMqClient),
			"SendAnakinHub": NewSendAnakinHub(rabbitMqClient),
		},
		rabbitMqClient: rabbitMqClient,
	}
}

//Todo: hubConfiguration must be dynamic and get from config db which run as sidecar
func (hub *Hub) ProcessByMessageType(messageContract *contracts.BridgeMessageContract) error {
	currentHub, exists := hub.hubs[messageContract.Type]
	if !exists {
		return errors.New("Hub could not be found for [Type]: " + messageContract.Type)
	}

	processResult := currentHub.Transmit(messageContract, &contracts.HubConfiguration{
		BrokerConfiguration: contracts.BrokerConfiguration{
			Nodes: []string{
				"localhost:5672",
			},
			ExchangeName: "noc-tools",
			ExchangeType: "topic",
			Durable:      true,
			AutoDelete:   false,
			Internal:     false,
			NoWait:       false,
			RoutingKey:   messageContract.RoutingKey,
			Mandatory:    false,
			Immediate:    false,
			ContentType:  "application/json",
		},
	})

	if processResult != nil {
		return processResult
	}

	return nil

}
