package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/satori/go.uuid"
	"net/http"
	"bytes"
	"io"
)

type Payload struct {
	Body string `json:"body"`
	Smid string `json:"smid"`
}

//go get -u gopkg.in/redis.v5
//go get -u github.com/labstack/echo
//go get -u github.com/satori/go.uuid
//go get -u github.com/streadway/amqp

func main() {
	server := echo.New()

	server.POST("/data", func(context echo.Context) error {
		payload := new(Payload)
		payloadUuid := uuid.NewV4().String()
		payloadString := bodyToString(context.Request().Body())
		if err := context.Bind(payload); err != nil {
			return err
		}

		go SaveToRedis(payloadUuid, payloadString)
		go PublishToRabbit(payloadString)

		return context.NoContent(http.StatusAccepted)
	})

	server.Run(fasthttp.New(":8080"))
}

func bodyToString(reader io.Reader) string {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)
	return buffer.String()
}