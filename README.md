# Overview Architecture
![Image alt text](./images/MessagingBridge.gif)

# grafana-webhook
grafana-webhook provides an easy way to write go http handlers for webhook channels

## usage example
Handle Grafana request and send a message by a service:

```go

...

// sending messages server
func main() {
	e := echo.New()
	bridgeService := services.NewMessageBridgeService()
	e.POST("/process", func(c echo.Context) error {
		messageContract := new(contracts.BridgeMessageContract)
		if err := c.Bind(messageContract); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err := bridgeService.Process(messageContract)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusAccepted, "")
	})

	go e.Logger.Fatal(e.Start(":1323"))
}

...

```

Above listener can be used to fill the url input in [Grafana Webhook notification channel](http://docs.grafana.org/alerting/notifications/#webhook).
Please look at [Message Bridge from Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/patterns/messaging/MessagingBridge.html).


## references
1. https://github.com/Azure/azure-event-hubs-go/blob/master/_examples/helloworld/readme.md
2. https://github.com/rabbitmq/rabbitmq-tutorials/blob/c665f54566903d21a4e716894f59dfb884adcb44/go/emit_log.go

## known issues
1. swagger
2. docker compose - rabbitmq,grafana,k8s scripts
3. config and logging

