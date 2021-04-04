package services

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"webhook/pkg/contracts"
	"webhook/pkg/hubs"
)

type MessageBridgeService struct {
	config *viper.Viper
	logger *zap.Logger
	hub    *hubs.Hub
}

func NewMessageBridgeService(config *viper.Viper, logger *zap.Logger) *MessageBridgeService {
	return &MessageBridgeService{
		hub: hubs.NewHub(config, logger),
	}
}

func (messageBridgeService *MessageBridgeService) Process(messageContract *contracts.BridgeMessageContract, correlationId string) error {
	messageContract.Type = getParseType()
	messageContract.RoutingKey = getRoutingKey()
	messageContract.CorrelationId = correlationId
	err := messageBridgeService.hub.ProcessByMessageType(messageContract, correlationId)
	if err != nil {
		return err
	}
	return nil
}

//todo: parse by log format
func getParseType() string {
	return "SendNetmonHub"
}

//todo: parse by log format and select right routing key
func getRoutingKey() string {
	return "netmon"
}
