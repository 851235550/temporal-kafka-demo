package producer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

const (
	kafkaTopic     = "temporal-topic"
	kafkaBrokerURL = "localhost:9092"
)

// ProduceMsg is the main workflow for producing messages to Kafka
func ProduceMsg(ctx context.Context) error {
	// Kafka writer configuration
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokerURL),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Start produce msg...")
	// Send 200 messages to Kafka
	for i := 0; i < 200; i++ {
		log.Println("Start produce msg-" + strconv.Itoa(i))
		message := kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", i)),
			Value: []byte(fmt.Sprintf("Message-%d", i)),
		}
		err := writer.WriteMessages(context.Background(), message)
		if err != nil {
			log.Printf("Error writing message to Kafka: %v", err)
			return err
		}
		log.Println("End produce msg-" + strconv.Itoa(i))
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
		return err
	}

	return nil
}

// CronProducerWorkflow is the cron workflow that triggers the producer activity every minute
func CronProducerWorkflow(ctx workflow.Context) error {
	// Define activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute the activity
	err := workflow.ExecuteActivity(ctx, ProduceMsg).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
