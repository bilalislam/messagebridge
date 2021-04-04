package hubs

import (
	"go.uber.org/zap"
	"time"
	"webhook/pkg/contracts"
	"webhook/pkg/infrastructure"
)

type SendNetmonHub struct {
	logger         *zap.Logger
	rabbitMqClient *infrastructure.RabbitMqClient
}

func NewSendNetmonHub(logger *zap.Logger, rabbitMqClient *infrastructure.RabbitMqClient) *SendNetmonHub {
	return &SendNetmonHub{
		logger:         logger,
		rabbitMqClient: rabbitMqClient,
	}
}

func (hub *SendNetmonHub) Transmit(message *contracts.BridgeMessageContract, hubConfiguration *contracts.HubConfiguration) error {
	return hub.rabbitMqClient.Publish(&hubConfiguration.BrokerConfiguration, &NetmonCommand{
		message.Message,
		message.Title,
		message.State,
		message.CorrelationId,
		time.Now(),
	}, message.CorrelationId)
}

type NetmonCommand struct {
	Message       string          `json:"message"`
	Title         string          `json:"title"`
	State         contracts.State `json:"state"`
	CorrelationId string          `json:"correlationId"`
	EventOn       time.Time       `json:"eventOn"`
}
