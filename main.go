package main

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"webhook/pkg/contracts"
	"webhook/services"
)

func main() {
	e := echo.New()
	bridgeService := services.NewMessageBridgeService()

	e.GET("/", Index)
	e.GET("/healthcheck", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
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
