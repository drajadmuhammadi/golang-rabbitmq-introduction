package services

import (
	"fmt"
	"golang-rabbitmq-introduction/consumer/consumer-topic-satu/config"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMq struct {
	base               *Base
	RabbitMqConnection *amqp.Connection
	RabbitMqChannel    *amqp.Channel
	RabbitQueue        amqp.Queue
	Delivery           <-chan amqp.Delivery
}

func (r *RabbitMq) InitializeRabbitMq() {
	var err error

	// Create Connection
	r.RabbitMqConnection, err = amqp.Dial("amqp://" +
		config.RABBITMQ_USER + ":" +
		config.RABBITMQ_PASSWORD + "@" +
		config.RABBITMQ_HOST + ":" +
		config.RABBITMQ_PORT + "/")
	if err != nil {
		log.Println(err)
	}
	// Open Channel
	r.openChannel()
	// Create Queue
	r.createQueue()
	// Bind Exchange to Queue
	r.bindExchangeToQueue()
}

func (r *RabbitMq) openChannel() {
	var err error
	fmt.Println("Topic name :", config.CONSUMER_TOPIC_NAME)
	fmt.Println("Queue name :", config.CONSUMER_ROUTING_KEY)
	fmt.Println("Handling routing key :", config.CONSUMER_QUEUE_NAME)
	r.RabbitMqChannel, err = r.RabbitMqConnection.Channel()
	if err != nil {
		log.Println(err)
	}
}

func (r *RabbitMq) createQueue() {
	var err error
	r.RabbitQueue, err = r.RabbitMqChannel.QueueDeclare(
		config.CONSUMER_QUEUE_NAME,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
}

func (r *RabbitMq) bindExchangeToQueue() {
	r.RabbitMqChannel.QueueBind(
		r.RabbitQueue.Name,
		config.CONSUMER_ROUTING_KEY,
		config.CONSUMER_TOPIC_NAME,
		false,
		nil,
	)
}

func (r *RabbitMq) StartConsume() {
	var err error
	r.Delivery, err = r.RabbitMqChannel.Consume(
		r.RabbitQueue.Name,
		config.CONSUMER_QUEUE_NAME,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	// Create Go Concurence to Listen Queue
	forever := make(chan bool)
	go func() {
		for d := range r.Delivery {
			r.base.Controller.Index(d, d.Body, r.RabbitMqChannel)
		}
	}()
	log.Println(" [*] Listening on topic :", config.CONSUMER_TOPIC_NAME, "with routing key :", config.CONSUMER_ROUTING_KEY)
	<-forever
}
