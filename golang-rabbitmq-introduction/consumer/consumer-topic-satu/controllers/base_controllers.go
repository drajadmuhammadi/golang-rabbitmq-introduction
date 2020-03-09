package controllers

import (
	"encoding/json"
	"golang-rabbitmq-introduction/consumer/consumer-topic-satu/models"
	"log"

	"github.com/streadway/amqp"
)

type BaseControllers struct {
	RabbitMqChannel *amqp.Channel
}

func (b *BaseControllers) Index(acknowledgement amqp.Delivery, reqData []byte, channel *amqp.Channel) {
	var methodAPI models.OpenMethodAPI
	var basicResponse models.BasicResponse
	// var checker bool
	json.Unmarshal(reqData, &methodAPI)
	b.RabbitMqChannel = channel

	switch methodAPI.Method {
	case "api_1":
		// valueData, _ := json.Marshal(methodAPI.Value)
		basicResponse.Status = true
		basicResponse.Data = methodAPI.Value
	default:
		basicResponse.Status = false
		basicResponse.Data = "Method Not Found"
	}

	acknowledgement.Ack(false)
	log.Println(acknowledgement.CorrelationId)
	// if sync, publish back to consumer
	if acknowledgement.Type == "sync" {
		responseConsumer, _ := json.Marshal(basicResponse)
		b.Publish("", "callback-mq-producer", responseConsumer, acknowledgement.CorrelationId)
	}
}

func (b *BaseControllers) Publish(topicName string, routingKey string, data []byte, corrId string) {
	b.RabbitMqChannel.Publish(
		topicName,  // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
			ReplyTo:     corrId,
		})
}
