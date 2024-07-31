package main

import (
	"log"
	"suger/config"
	"suger/producer"
)

func main() {
	// Initialize Temporal client
	c, err := config.NewTemporalClient()
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	log.Println("Start work...")

	// go producer.ProducerWorkflow()

	// go consumer.ConsumerWorkflow()

	// Start producer and consumer workers
	go producer.StartWorker(c)
	// go consumer.StartWorker(client)

	// Block the main thread
	select {}
}
