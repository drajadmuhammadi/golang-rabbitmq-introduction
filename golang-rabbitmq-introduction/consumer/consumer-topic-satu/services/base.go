package services

import "golang-rabbitmq-introduction/consumer/consumer-topic-satu/controllers"

type Base struct {
	Controller controllers.BaseControllers
	RabbitMq   RabbitMq
}

func (b *Base) RunMainService() {
	b.RabbitMq.InitializeRabbitMq()

	b.registerBaseService()
	b.RabbitMq.StartConsume()
}

func (base *Base) registerBaseService() {
	base.RabbitMq.base = base
}
