package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/workflow"
)

const (
	kafkaTopic     = "temporal-topic"
	kafkaBrokerURL = "localhost:9092"
)

// ProduceMsg is the main workflow for producing messages to Kafka
func ProduceMsg(_ context.Context, msg string) error {
	// fmt.Println("Producing message:", msg)
	// Kafka writer configuration
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokerURL),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	log.Printf("Start produce msg: %s\n", msg)
	message := kafka.Message{
		Value: []byte(msg),
	}
	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Error writing message to Kafka: %v", err)
		return err
	}
	log.Printf("End produce msg: %s\n", msg)

	// Close the writer
	if err := writer.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
		return err
	}

	return nil
}

func CronParentProducerWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	childWorkflowResults := make([]workflow.Future, 0, 200)
	// Start 200 child workflows asynchronously
	for i := 0; i < 200; i++ {
		// Define child workflow options
		childWorkflowOptions := workflow.ChildWorkflowOptions{
			WorkflowID: fmt.Sprintf("child-%d", i),
			TaskQueue:  "child-task-queue",
		}
		childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
		future := workflow.ExecuteChildWorkflow(childCtx, ChildWorkflow, fmt.Sprintf("msg-%d", i))
		childWorkflowResults = append(childWorkflowResults, future)
	}

	// Wait for all child workflows to complete
	for i, future := range childWorkflowResults {
		err := future.Get(ctx, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute child workflow %d: %v", i, err))
			return err
		}
	}

	return nil
}

func ChildWorkflow(ctx workflow.Context, msg string) error {
	// Define activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute the activity
	err := workflow.ExecuteActivity(ctx, ProduceMsg, msg).Get(ctx, nil)
	if err != nil {
		logger := workflow.GetLogger(ctx)
		logger.Error("Failed to execute activity ProduceMsg", "Error", err)
	}
	return err
}

// // CronProducerWorkflow is the cron workflow that triggers the producer activity every minute
// func CronProducerWorkflow(ctx workflow.Context) error {
// 	// Define activity options
// 	ao := workflow.ActivityOptions{
// 		StartToCloseTimeout: time.Second * 20,
// 	}
// 	ctx = workflow.WithActivityOptions(ctx, ao)

// 	// Execute the activity
// 	return workflow.ExecuteActivity(ctx, ProduceMsg).Get(ctx, nil)
// }
