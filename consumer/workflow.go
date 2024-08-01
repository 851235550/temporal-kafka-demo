package consumer

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

func ConsumerMsg(_ context.Context) error {
	// Kafka reader configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokerURL},
		Topic:   kafkaTopic,
	})

	msg, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Error reading message from Kafka: %v", err)
		return err
	}
	log.Printf("Message received: key=%s value=%s offset=%d", string(msg.Key), string(msg.Value), msg.Offset)

	// Close the reader
	if err := reader.Close(); err != nil {
		log.Printf("Error closing Kafka reader: %v", err)
		return err
	}

	return nil
}

func CronParentConsumerWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	childWorkflowResults := make([]workflow.Future, 0, 200)
	for i := 0; i < 200; i++ {
		opts := workflow.ChildWorkflowOptions{
			WorkflowID: fmt.Sprintf("consumer-child-%d", i),
			TaskQueue:  "consumer-child-task-queue",
		}
		childCtx := workflow.WithChildOptions(ctx, opts)
		future := workflow.ExecuteChildWorkflow(childCtx, ChildWorkflow)
		childWorkflowResults = append(childWorkflowResults, future)
	}

	// Wait for all child workflows to complete
	for i, future := range childWorkflowResults {
		err := future.Get(ctx, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute consumer child workflow %d: %v", i, err))
			return err
		}
	}

	return nil
}

func ChildWorkflow(ctx workflow.Context) error {
	// Define activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute the activity
	err := workflow.ExecuteActivity(ctx, ConsumerMsg).Get(ctx, nil)
	if err != nil {
		logger := workflow.GetLogger(ctx)
		logger.Error("Failed to execute activity ConsumerMsg", "Error", err)
	}

	return err
}

// CronConsumerWorkflow is the cron workflow that triggers the consumer workflow every two minutes
// func CronConsumerWorkflow(ctx workflow.Context) error {
// 	cronOpts := workflow.ChildWorkflowOptions{
// 		CronSchedule: "*/2 * * * *", // Every two minutes
// 	}
// 	ctx = workflow.WithChildOptions(ctx, cronOpts)
// 	return workflow.ExecuteChildWorkflow(ctx, ConsumerWorkflow).Get(ctx, nil)
// }
