package consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

const (
	kafkaTopic     = "temporal-topic"
	kafkaBrokerURL = "localhost:9092"
)

// ConsumerWorkflow is the main workflow for consuming messages from Kafka
func ConsumerWorkflow(ctx workflow.Context) error {
	// Kafka reader configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokerURL},
		Topic:   kafkaTopic,
	})

	// Read 200 messages from Kafka
	for i := 0; i < 200000; i++ {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			return err
		}
		log.Printf("Message received: key=%s value=%s", string(msg.Key), string(msg.Value))
	}

	// Close the reader
	if err := reader.Close(); err != nil {
		log.Printf("Error closing Kafka reader: %v", err)
		return err
	}

	return nil
}

// CronConsumerWorkflow is the cron workflow that triggers the consumer workflow every two minutes
func CronConsumerWorkflow(ctx workflow.Context) error {
	cronOpts := workflow.ChildWorkflowOptions{
		CronSchedule: "*/2 * * * *", // Every two minutes
	}
	ctx = workflow.WithChildOptions(ctx, cronOpts)
	return workflow.ExecuteChildWorkflow(ctx, ConsumerWorkflow).Get(ctx, nil)
}
