package main

import (
	"log"
	"os"
	"os/signal"
	"suger/config"
	"suger/consumer"
	"suger/producer"
	"syscall"
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
	go consumer.StartWorker(c)

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block the main thread until a signal is received
	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down...", sig)

	// Perform any cleanup or graceful shutdown here
	c.Close()
	log.Println("Client closed. Exiting...")
}
