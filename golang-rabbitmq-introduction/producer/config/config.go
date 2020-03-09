package config

// API Endpoint Config
const API_ENDPOINT_HTTP_HOST = "127.0.0.1"
const API_ENDPOINT_HTTP_PORT = "8080"

// RabbitMQ Config
const RABBITMQ_HOST = "192.168.141.136"
const RABBITMQ_PORT = "5672"
const RABBITMQ_USER = "mqadmin"
const RABBITMQ_PASSWORD = "mqadmin"

// List Topic
func ListTopicName() []string {
	var ListTopicName []string

	ListTopicName = append(ListTopicName, "topic-satu")
	ListTopicName = append(ListTopicName, "topic-dua")

	return ListTopicName
}
