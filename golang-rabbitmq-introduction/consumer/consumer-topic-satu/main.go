package main

import (
	"golang-rabbitmq-introduction/consumer/consumer-topic-satu/services"
	"log"
)

func main() {
	log.Println("Starting Consumer Topic Satu ...")
	mainServices := services.Base{}
	mainServices.RunMainService()
}
