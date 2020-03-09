package validator

import (
	"golang-rabbitmq-introduction/producer/models"
)

func BasicValidator(dataAPI models.OpenRequestAPI) (bool, string) {
	if len(dataAPI.TopicName) == 0 {
		return false, "'topic_name' field is empty"
	}
	if len(dataAPI.RoutingKey) == 0 {
		return false, "'routing_key' field is empty"
	}
	return topicNameValidator(dataAPI.TopicName, dataAPI.Data)
}

func topicNameValidator(topicName string, data interface{}) (bool, string) {
	switch topicName {
	case "topic-satu":
		return topicSatuValidator(data)
	case "topic-dua":
		return true, ""
	default:
		return false, "Invalid Topic Name, Please Contact Administrator"
	}
}
