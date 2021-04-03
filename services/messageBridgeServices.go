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
	err := messageBridgeService.hub.ProcessByMessageType(messageContract)
	if err != nil {
		return err
	}
	return nil
}
