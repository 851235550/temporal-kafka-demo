package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"suger/config"
	"suger/consumer"
	"suger/producer"
	"syscall"
)

var workflowCnt int

func main() {
	flag.IntVar(&workflowCnt, "wf-cnt", 200, "The number of sub-producer and sub-consumer workflows. The default value is 200, which means there will be 200 sub-production flows and 200 sub-consumption flows")
	flag.Parse()

	consumer.SetChildWorkerCnt(workflowCnt)
	producer.SetChildWorkerCnt(workflowCnt)

	// Initialize Temporal client
	c, err := config.NewTemporalClient()
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	log.Println("Start work...")

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
