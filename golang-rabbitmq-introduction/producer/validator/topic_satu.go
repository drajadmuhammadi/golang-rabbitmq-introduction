package validator

import (
	"encoding/json"
	"golang-rabbitmq-introduction/producer/models"
	"reflect"
)

func topicSatuValidator(data interface{}) (bool, string) {
	var openData models.OpenTopicSatu
	byteData, _ := json.Marshal(data)
	json.Unmarshal(byteData, &openData)

	if len(openData.Method) == 0 {
		return false, "'method' field on data Topic Satu is empty"
	}

	if reflect.ValueOf(openData.Value).IsNil() {
		return false, "'value' field on data Topic Satu is empty"
	}

	return true, ""
}
