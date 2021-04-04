package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"time"
	"webhook/pkg/contracts"
)

type RabbitMqClient struct {
	logger *zap.Logger
}

func NewRabbitMqClient(logger *zap.Logger) *RabbitMqClient {
	return &RabbitMqClient{
		logger: logger,
	}
}

func (rabbitMqClient *RabbitMqClient) Publish(
	brokerConfiguration *contracts.BrokerConfiguration,
	message interface{},
	correlationId string) error {

	environment := os.Getenv("ENV_FILE")
	if len(environment) == 0 || environment == "dev" {
		environment = "development"
	}

	routingKey := brokerConfiguration.RoutingKey + "-" + environment

	rabbitMqClient.logger.Info(fmt.Sprintf("[RABBITMQ] connection trying to create to bus [Broker-RoutingKey]: %s", routingKey), zap.String("fields.CorrelationId", correlationId))

	node := chooseRandomNode(brokerConfiguration.Nodes)
	connection, err := amqp.Dial(node)
	if err != nil {
		rabbitMqClient.logger.Error(fmt.Sprintf("[RABBITMQ] could not be create a connection to bus [Error]: %s", err.Error()), zap.String("fields.CorrelationId", correlationId))

		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		rabbitMqClient.logger.Error(fmt.Sprintf("[RABBITMQ] could not be create a channel in bus  [Error]: %s", err.Error()), zap.String("fields.CorrelationId", correlationId))
		return err
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		brokerConfiguration.ExchangeName,
		brokerConfiguration.ExchangeType,
		brokerConfiguration.Durable,
		brokerConfiguration.AutoDelete,
		brokerConfiguration.Internal,
		brokerConfiguration.NoWait,
		nil,
	)
	if err != nil {
		rabbitMqClient.logger.Error(fmt.Sprintf("[RABBITMQ] could not decleare an exchange in bus  [Error]: %s", err.Error()), zap.String("fields.CorrelationId", correlationId))
		return err
	}

	if err := channel.Confirm(false); err != nil {
		rabbitMqClient.logger.Error(fmt.Sprintf("[RABBITMQ] channel could not be put into confirm mode in bus  [Error]: %s", err.Error()), zap.String("fields.CorrelationId", correlationId))
	}

	data, err := json.Marshal(message)
	if err != nil {
		rabbitMqClient.logger.Error(fmt.Sprintf("[RABBITMQ] message could not be converted to verified format <tips: []byte> [Error]: %s", err.Error()), zap.String("fields.CorrelationId", correlationId))
		return err
	}
	err = channel.Publish(
		brokerConfiguration.ExchangeName,
		routingKey,
		brokerConfiguration.Mandatory,
		brokerConfiguration.Immediate,
		amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  brokerConfiguration.ContentType,
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		rabbitMqClient.logger.Error(fmt.Sprintf("[RABBITMQ] message could not published to bus  [Error]: %s", err.Error()), zap.String("fields.CorrelationId", correlationId))
		return err
	}

	rabbitMqClient.logger.Info(fmt.Sprint("[RABBITMQ] message published to bus "), zap.String("fields.CorrelationId", correlationId))
	return nil
}

func chooseRandomNode(nodes []string) string {
	//defaultNode := nodes[0]
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })

	//if defaultNode == nodes[0] {
	//	return nodes[1]
	//}

	return nodes[0]
}
