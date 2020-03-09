package services

import (
	"golang-rabbitmq-introduction/producer/config"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMq struct {
	base               *Base
	RabbitMqConnection *amqp.Connection
	RabbitMqChannel    *amqp.Channel
	Callback           <-chan amqp.Delivery
	ActiveTopic        []string
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
	r.OpenChannel()

	// Initiate Consumer to handle callback consumer
	r.registerCallback()

	// Create Topic
	r.RegisterDefaultExchange()

}

func (r *RabbitMq) OpenChannel() {
	var err error
	r.RabbitMqChannel, err = r.RabbitMqConnection.Channel()
	if err != nil {
		log.Println(err)
	}
}

func (r *RabbitMq) RegisterDefaultExchange() {
	ListTopicName := config.ListTopicName()
	for _, valueListTopicName := range ListTopicName {
		isExchangeActivated := r.CheckAndActivateExchange(valueListTopicName)
		if !isExchangeActivated {
			r.RegisterTopicName(valueListTopicName)
		}
	}
}

func (r *RabbitMq) CheckAndActivateExchange(exchangeName string) bool {
	isExist := r.RabbitMqChannel.ExchangeDeclarePassive(
		exchangeName, // exchangeName
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	// Jika ada cukup aktivasi saja
	if isExist == nil {
		return true
	} else {
		r.OpenChannel()
		return false
	}
}

func (r *RabbitMq) RegisterTopicName(topicName string) {
	if !r.isTopicCreated(topicName) {
		r.ActiveTopic = append(r.ActiveTopic, topicName)
		r.RabbitMqChannel.ExchangeDeclare(
			topicName, // name
			"topic",   // type
			true,      // durable
			false,     // auto-deleted
			false,     // internal
			false,     // no-wait
			nil,       // arguments
		)
		log.Println("Initiate new topic", topicName)
	}

}

func (r *RabbitMq) isTopicCreated(topicName string) bool {
	for _, activeTopicName := range r.ActiveTopic {
		if activeTopicName == topicName {
			return true
		}
	}
	return false
}

// For Publishing to Queue
func (r *RabbitMq) Publish(topicName string, routingKey string, data []byte, corrId string, isSync string) {
	abc := r.RabbitMqChannel.Publish(
		topicName,  // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Body:          data,
			ContentType:   "text/plain",
			CorrelationId: corrId,
			// deliver mode 1 : temporary, 2 : permanent
			DeliveryMode: 2,
			ReplyTo:      corrId,
			Type:         isSync,
		})
	log.Println(abc)
}

// Consume Callback from Consumer
func (r *RabbitMq) registerCallback() {
	r.RabbitMqChannel.QueueDeclare(
		"callback-mq-producer",
		true,
		false,
		false,
		false,
		nil,
	)
	r.StartConsumeCallback()
	log.Println("Create Queue callback")
}

func (r *RabbitMq) StartConsumeCallback() {
	r.Callback, _ = r.RabbitMqChannel.Consume(
		"callback-mq-producer",
		"callback-consumer",
		true,
		true,
		false,
		false,
		nil,
	)
}
