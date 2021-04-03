package services

import (
	"webhook/pkg/contracts"
	"webhook/pkg/hubs"
)

type MessageBridgeService struct {
	hub *hubs.Hub
}

func NewMessageBridgeService() *MessageBridgeService {
	return &MessageBridgeService{
		hub: hubs.NewHub(),
	}
}

func (messageBridgeService *MessageBridgeService) Process(messageContract *contracts.BridgeMessageContract) error {
	messageContract.Type = getParseType()
	messageContract.RoutingKey = getRoutingKey()
	err := messageBridgeService.hub.ProcessByMessageType(messageContract)
	if err != nil {
		return err
	}
	return nil
}

func getParseType() string {
	return "SendNetmonHub"
}

func getRoutingKey() string {
	return "netmon"
}
