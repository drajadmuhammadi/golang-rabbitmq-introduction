package main

import (
	"golang-rabbitmq-introduction/producer/services"
	"log"
)

func main() {
	// Start Main Services
	log.Println("Starting Server Produser...")
	mainServices := services.Base{}
	mainServices.RunMainService()
}
