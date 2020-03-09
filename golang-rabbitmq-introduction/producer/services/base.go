package services

import (
	"golang-rabbitmq-introduction/producer/config"
	"net/http"
)

type Base struct {
	RabbitMq   RabbitMq
	HttpServer HttpServer
}

func (b *Base) RunMainService() {
	b.RabbitMq.InitializeRabbitMq()

	b.registerBaseService()
	b.HttpServer.configureHandler()

	http.ListenAndServe(config.API_ENDPOINT_HTTP_HOST+":"+config.API_ENDPOINT_HTTP_PORT, nil)
}

func (b *Base) registerBaseService() {
	b.RabbitMq.base = b
	b.HttpServer.base = b
}
