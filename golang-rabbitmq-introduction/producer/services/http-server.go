package services

import (
	"encoding/json"
	"fmt"
	"golang-rabbitmq-introduction/producer/models"
	"golang-rabbitmq-introduction/producer/validator"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type HttpServer struct {
	base *Base
}

func (h *HttpServer) configureHandler() {
	http.HandleFunc("/", h.rootHandler)
}

func (h *HttpServer) rootHandler(w http.ResponseWriter, r *http.Request) {

	var response models.BasicResponse
	var jsonResponse []byte

	reqBody, _ := ioutil.ReadAll(r.Body)

	// Open Request From API
	openRequestAPI := models.OpenRequestAPIDefault()
	json.Unmarshal(reqBody, &openRequestAPI)

	// Check Validation Request API
	isSuccess, errValidator := validator.BasicValidator(openRequestAPI)
	corrId := h.randomString(16)
	log.Println(corrId)
	if isSuccess {
		dataByte, _ := json.Marshal(openRequestAPI.Data)
		// Publish Request to Rabbit MQ
		h.base.RabbitMq.Publish(
			openRequestAPI.TopicName,  // Topic Name
			openRequestAPI.RoutingKey, // Routing Key
			dataByte,                  // Data <Byte>
			corrId,                    // Correlation ID
			openRequestAPI.Type,       // type sync or async
		)

		// Default Message For async
		response.Status = isSuccess
		response.Value = "Success receive data"
	} else {
		response.Status = isSuccess
		response.Value = errValidator
		fmt.Println("Error Message :", string(errValidator))
	}

	if openRequestAPI.Type == "async" {
		jsonResponse, _ = json.Marshal(response)
	} else {
		json.Unmarshal(h.base.HttpServer.consumeCallback(corrId), &response.Value)
		jsonResponse, _ = json.Marshal(response)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Func to Generate Correlation ID
func (h *HttpServer) randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(h.randInt(65, 90))
	}
	return string(bytes)
}

func (h *HttpServer) randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// Consume Callback Sync
func (h *HttpServer) consumeCallback(corrId string) []byte {
	var callbackResponse []byte

	callbackChannel := make(chan bool)
	go func() {
		for c := range h.base.RabbitMq.Callback {
			if c.ReplyTo == corrId {
				callbackResponse = c.Body
				close(callbackChannel)
				c.Ack(false)
				break
			} else {
				c.Nack(false, true)
			}
		}
	}()
	<-callbackChannel

	return callbackResponse
}
