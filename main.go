package main

import (
	"github.com/bilalislam/torc/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"webhook/configuration"
	"webhook/pkg/contracts"
	"webhook/pkg/services"
)

func main() {
	logger := log.GetLogger()
	configClient := configuration.NewConfigClient(".env", ".")
	h := configuration.ConfigHandler{
		ConfigClient: configClient,
	}
	config, err := h.NewConfigHandler()
	if err != nil {
		panic(err.Error())

	}
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"@timestamp":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"fields.CorrelationId":"${header:x-correlation-id}","error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))

	bridgeService := services.NewMessageBridgeService(config, logger)

	e.Use(middleware.Recover())
	e.GET("/", Index)
	e.GET("/healthcheck", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/process", func(c echo.Context) error {
		messageContract := new(contracts.BridgeMessageContract)
		if err := c.Bind(messageContract); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		c.Request().Header.Set("x-correlation-id", random.String(32))
		correlationId := c.Request().Header.Get("x-correlation-id")
		err := bridgeService.Process(messageContract, correlationId)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusAccepted, "")
	})

	env := os.Getenv("ENV_FILE")
	var port string
	if env == "" || env == "dev" {
		port = "8080"
	} else {
		port = "80"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func Index(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, `{"status":"healthy !"}`)
}
