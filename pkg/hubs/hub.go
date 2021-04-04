package hubs

import (
	"errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"webhook/pkg/contracts"
	"webhook/pkg/infrastructure"
)

type IHub interface {
	Transmit(message *contracts.BridgeMessageContract, hubConfiguration *contracts.HubConfiguration) error
}

type Hub struct {
	config         *viper.Viper
	logger         *zap.Logger
	hubs           map[string]IHub
	rabbitMqClient *infrastructure.RabbitMqClient
}

func NewHub(config *viper.Viper, logger *zap.Logger) *Hub {
	rabbitMqClient := infrastructure.NewRabbitMqClient(logger)
	return &Hub{
		config: config,
		hubs: map[string]IHub{
			"SendNetmonHub": NewSendNetmonHub(logger, rabbitMqClient),
			"SendAnakinHub": NewSendAnakinHub(logger, rabbitMqClient),
		},
		rabbitMqClient: rabbitMqClient,
	}
}

//Todo: hubConfiguration must be dynamic and get from config db which run as sidecar
func (hub *Hub) ProcessByMessageType(messageContract *contracts.BridgeMessageContract, correlationId string) error {
	currentHub, exists := hub.hubs[messageContract.Type]
	if !exists {
		hub.logger.Error("[HUB] could not be found for [Type]: "+messageContract.Type, zap.String("fields.CorrelationId", correlationId))
		return errors.New("Hub could not be found for [Type]: " + messageContract.Type)
	}

	processResult := currentHub.Transmit(messageContract, &contracts.HubConfiguration{
		BrokerConfiguration: contracts.BrokerConfiguration{
			RoutingKey:   messageContract.RoutingKey,
			Nodes:        hub.config.GetStringSlice("rabbitmq.nodes"),
			ExchangeName: hub.config.GetString("rabbitmq.exchange-name"),
			ExchangeType: hub.config.GetString("rabbitmq.exchange-type"),
			ContentType:  hub.config.GetString("rabbitmq.content-type"),
			Durable:      hub.config.GetBool("rabbitmq.durable"),
			AutoDelete:   hub.config.GetBool("rabbitmq.auto-delete"),
			Internal:     hub.config.GetBool("rabbitmq.internal"),
			NoWait:       hub.config.GetBool("rabbitmq.no-wait"),
			Mandatory:    hub.config.GetBool("rabbitmq.mandatory"),
			Immediate:    hub.config.GetBool("rabbitmq.immediate"),
		},
	})

	if processResult != nil {
		return processResult
	}

	return nil

}
