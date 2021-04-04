package hubs

import (
	"go.uber.org/zap"
	"time"
	"webhook/pkg/contracts"
	"webhook/pkg/infrastructure"
)

type SendAnakinHub struct {
	logger         *zap.Logger
	rabbitMqClient *infrastructure.RabbitMqClient
}

func NewSendAnakinHub(logger *zap.Logger, rabbitMqClient *infrastructure.RabbitMqClient) *SendAnakinHub {
	return &SendAnakinHub{
		logger:         logger,
		rabbitMqClient: rabbitMqClient,
	}
}

func (hub *SendAnakinHub) Transmit(message *contracts.BridgeMessageContract, hubConfiguration *contracts.HubConfiguration) error {
	return hub.rabbitMqClient.Publish(&hubConfiguration.BrokerConfiguration, &AnakinCommand{
		message.Message,
		message.Title,
		message.State,
		message.CorrelationId,
		time.Now(),
	}, message.CorrelationId)
}

type AnakinCommand struct {
	Message       string          `json:"message"`
	Title         string          `json:"title"`
	State         contracts.State `json:"state"`
	CorrelationId string          `json:"correlationId"`
	EventOn       time.Time       `json:"eventOn"`
}
